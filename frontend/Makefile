.PHONY: init
init:
	npm ci

public/bundle.%: $(shell find src)
	npm run build

.PHONY: deploy
deploy: public/bundle.js public/bundle.css public/index.html
	aws s3 sync public s3://sthub.in.fkn.space/
	aws cloudfront create-invalidation --distribution-id E2OHW6JPREC4E4 --paths /index.html /bundle.css /bundle.js

.PHONY: build
build: public/bundle.js public/bundle.css public/index.html

all: deploy
