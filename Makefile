.PHONY: init
init:
	go get github.com/GeertJohan/go.rice/rice
	go get github.com/akavel/rsrc

.PHONY: build
build:
	# rice embed-go
	GOOS=windows GOARCH=amd64 go build -ldflags -H=windowsgui

.PHONY: icon
icon:
	rsrc -ico=assets/sthub.ico