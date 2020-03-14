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

.PHONY: get-testships
get-testships:
	@curl -s -X POST -d "$(shell cat ./current_testships.txt)" https://bgaaa1fe37.execute-api.eu-central-1.amazonaws.com/prod/update-testships
