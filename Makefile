build:
	go build

release:
	CGO_ENABLED=0 go build -a -ldflags='-extldflags "-static" -s -w' -installsuffix netgo

build_docker: release
	docker build -t "toybox" .

busybox_test:
	cd testdata && git clone git://busybox.net/busybox.git
	find testdata/busybox/testsuite -type f | xargs -I{} sed -i "s/busybox/toybox/g" {}
