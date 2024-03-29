# Description: Makefile for gosh
# Author: Andrew Leung
# Maintainer: Carlos Morales

## @DIR_BUILD Output directory
DIR_BUILD= ./builds

## @DIR_PKG_PLUGIN so files output directory
DIR_PKG_PLUGIN=./pkg/plugins

## @DIR_CONF_LINUX Local user configuration directory Linux
DIR_CONF_LINUX=${HOME}/.config/gosh
## @DIR_PLUGIN_LINUX Local user plugin directory Linux
DIR_PLUGIN_LINUX=${HOME}/.config/gosh/plugins

GOPATH := $(shell go env GOPATH)
APP=gosh
MAIN=cmd/gosh.go

.PHONY: build-all run plugin
all: clean build-all plugin

clean:
	rm -rf $(DIR_PKG)
	rm -rf $(DIR_BUILD)

build-all: build-linux
build-linux:
	mkdir -p $(DIR_BUILD)/linux
	GOOS=linux go build -o $(DIR_BUILD)/linux/$(APP) $(MAIN)

plugin:
	mkdir -p $(DIR_PKG_PLUGIN)
	./scripts/make.sh
	mv $(DIR_PKG_PLUGIN) $(DIR_BUILD)/

install:
	mkdir -p $(DIR_PLUGIN)
	cp $(DIR_PKG)/plugins/* $(DIR_PLUGIN_LINUX)/
	go install $(MAIN)

uninstall:
	rm -rf $(DIR_PLUGIN_LINUX)
	rm ${GOPATH}/bin/gosh

