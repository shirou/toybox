VERSION=0.0.1

.PHONY: build
build:
	go build

.PHONY: test
test:
	go test ./...

init:
	go get github.com/Songmu/goxz/cmd/goxz

.PHONE: release
release: init
	CGO_ENABLED=0 goxz -pv=$(VERSION) -os=freebsd,darwin,linux -arch=amd64 -d=dist -build-ldflags '-extldflags "-static" -s -w'

.PHONE: build_docker
build_docker:
	docker build -t "toybox" .
