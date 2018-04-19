.PHONY: build
build:
	go build

.PHONY: test
test:
	go test ./...

.PHONE: release
release:
	CGO_ENABLED=0 go build -a -ldflags='-extldflags "-static" -s -w' -installsuffix netgo

.PHONE: build_docker
build_docker: release
	docker build -t "toybox" .
