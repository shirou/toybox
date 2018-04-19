VERSION=0.0.1

.PHONY: build
build:
	go build

.PHONY: test
test:
	go test ./...

init:
	go get github.com/Songmu/goxz/cmd/goxz
	go get github.com/tcnksm/ghr

.PHONE: release
release: init
	/bin/rm -rf dist
	CGO_ENABLED=0 goxz -pv=$(VERSION) -os=freebsd,darwin,linux -arch=amd64 -d=dist -build-ldflags '-extldflags "-static" -s -w'
	source .env
	ghr $(VERSION) ./dist

.PHONE: build_docker
build_docker:
	docker build -t "toybox" .
