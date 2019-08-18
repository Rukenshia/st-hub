.PHONY: seed
seed:
	./seed.sh

.PHONY: build
build:
	rice embed-go
	go build