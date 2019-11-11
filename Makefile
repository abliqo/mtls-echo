BINARY_PATH=bin/mtls-echo

.PHONY: all
all: clean test build

.PHONY: clean
clean:
	go clean
	rm -f $(BINARY_PATH)

.PHONY: test
test:
	go test -v ./...

build: test $(BINARY_PATH)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BINARY_PATH) -v

.PHONY: run
run: build
	bin/mtls-echo

.PHONY: docker-build
docker-build: build
	docker build -t mtls-echo ./

.PHONY: docker-run
docker-run:
	./docker-run.sh
