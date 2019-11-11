package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"
)

// Cert captures a subset of ready-to-render fields from x509.Certificate
type Cert struct {
	Version            string
	SerialNumber       string
	Subject            string
	Issuer             string
	NotBefore          string
	NotAfter           string
	PublicKeyAlgorithm string
	PublicKeyPem       string
}

func FromX509Certificate(c *x509.Certificate) *Cert {
	pubkeyBytes, _ := x509.MarshalPKIXPublicKey(c.PublicKey)
	p := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubkeyBytes,
	}

	res := &Cert{
		Version:            fmt.Sprintf("%v", c.Version),
		SerialNumber:       fmt.Sprintf("%v", c.SerialNumber),
		Subject:            c.Subject.String(),
		Issuer:             c.Issuer.String(),
		NotBefore:          c.NotBefore.Format(time.RFC1123),
		NotAfter:           c.NotAfter.Format(time.RFC1123),
		PublicKeyAlgorithm: c.PublicKeyAlgorithm.String(),
		PublicKeyPem:       string(pem.EncodeToMemory(p)),
	}
	return res
}
