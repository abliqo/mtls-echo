package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/spf13/viper"
)

const CONFIG_VAR_NAME = "CONFIG_PATH"
const DEFAULT_CONFIG_PATH = "config/dev.json"

func main() {
	loadConfig()

	log.Println("Started")
	go startRedirectServer()
	go startTLSServer()

	quitChannel := make(chan os.Signal)
	signal.Notify(quitChannel, os.Interrupt, os.Kill)

	// Wait for OS interrupt signal
	<-quitChannel
	log.Println("Stopped")
}

func loadConfig() {
	configPath := os.Getenv(CONFIG_VAR_NAME)
	if len(configPath) == 0 {
		configPath = DEFAULT_CONFIG_PATH
	}

	viper.SetConfigFile(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Failed to read configuration: %v", err)
	}
}

func startRedirectServer() {
	redirectAddress := fmt.Sprintf("%s:%d", viper.GetString("server.host"), viper.GetInt("server.https_port"))
	redir := http.NewServeMux()
	redir.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		url := fmt.Sprintf("https://%s", redirectAddress)
		http.Redirect(w, r, url, http.StatusMovedPermanently)
	})

	server := http.Server{
		Addr:    fmt.Sprintf("%s:%d", viper.GetString("server.host"), viper.GetInt("server.http_port")),
		Handler: redir,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to start redirect server: %v", err)
	}
}

func startTLSServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/json", jsonHandler)

	server := http.Server{
		Addr: fmt.Sprintf("%s:%d", viper.GetString("server.host"), viper.GetInt("server.https_port")),
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			ClientAuth: tls.RequestClientCert,
		},
		Handler: mux,
	}

	err := server.ListenAndServeTLS(
		viper.GetString("server.cert"),
		viper.GetString("server.key"))
	if err != nil {
		log.Fatalf("Failed to start TLS server: %v", err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if len(r.TLS.PeerCertificates) == 0 {
		w.WriteHeader(http.StatusPreconditionRequired)
		content, err := ioutil.ReadFile("content/nocert.html")
		if err != nil {
			log.Fatalf("Failed to load content: %v", err)
		}
		w.Write([]byte(content))
		return
	}

	c := FromX509Certificate(r.TLS.PeerCertificates[0])
	t, err := template.ParseFiles("content/success.tpl.html")
	if err != nil {
		log.Fatalf("Failed to load template: %v", err)
	}

	p := Page{
		Title:        "Mutual TLS Echo",
		Cert:         c,
		JsonEndpoint: fmt.Sprintf("https://%s:%d/json", viper.GetString("server.host"), viper.GetInt("server.https_port")),
	}
	err = t.Execute(w, p)
	if err != nil {
		log.Fatalf("Failed to render template: %v", err)
	}
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	setCommonHeaders(w)
	w.Header().Set("content-type", "application/json")

	if len(r.TLS.PeerCertificates) == 0 {
		w.WriteHeader(http.StatusPreconditionRequired)
		content := `{
	"error": "Client certificate not provided.",
	"help": "https://github.com/abliqo/mtls-echo/user-guide.md"
}`
		w.Write([]byte(content))
		return
	}

	c := FromX509Certificate(r.TLS.PeerCertificates[0])
	jsonBytes, err := json.Marshal(c)
	if err != nil {
		log.Fatalf("Failed to convert certificate to json: %v", err)
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonBytes)
	if err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func setCommonHeaders(w http.ResponseWriter) {
	w.Header().Set("cache-control", "no-cache")
	w.Header().Set("server", "mtls-echo")
}
