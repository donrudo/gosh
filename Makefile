DIR_CONF=${HOME}/.config/gosh
DIR_PLUGIN=${HOME}/.config/gosh/plugins
GOPATH := $(shell go env GOPATH)
DIR_PKG=builds

APP=gosh
MAIN=cmd/gosh.go

.PHONY: build-all run plugin
all: clean build-all plugin

clean:
	rm -rf $(DIR_PKG)

build-all: build-linux
build-linux:
	GOOS=linux go build -o $(DIR_PKG)/linux/$(APP) $(MAIN)

plugin:
	./scripts/make.sh

install:
	mkdir -p $(DIR_PLUGIN)
	cp $(DIR_PKG)/plugins/* $(DIR_PLUGIN)/
	go install $(MAIN)

uninstall:
	rm -rf $(DIR_PLUGIN)
	rm ${GOPATH}/bin/gosh

