.PHONY: all run vet test build cmd tools generate

all: vet build

run: build
	bin/gwemu --cert-file=clientcert.crt --chain=chain.crt --key-file=private.key --span-endpoint=localhost:64964

vet:
	go vet ./...
	revive ./...
	golint ./...

test:
	go test -timeout 10s ./...

build: cmd

cmd:
	cd cmd/gwemu && go build -o ../../bin/gwemu

tools: 
	cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go get %

generate:
	buf lint
	buf generate