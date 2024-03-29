# Description: Makefile for gosh
# Author: Andrew Leung
# Maintainer: Carlos Morales

## @DIR_BUILD Output directory
DIR_BUILD = ./builds
DIR_BUILD_PLUGIN = ./builds/plugins

## @DIR_PKG_PLUGIN so files output directory
DIR_PKG_PLUGIN = ./pkg/plugins
DIR_SRC_PLUGIN = ./plugins
SRC_PLUGIN := $(wildcard $(DIR_SRC_PLUGIN)/*.go)
OBJ_PLUGIN := $(SRC_PLUGIN:$(DIR_SRC_PLUGIN)/%.go=%)
# OBJ_PLUGIN := $(patsubst $(DIR_SRC_PLUGIN)/%.go, $(DIR_PKG_PLUGIN)/%.so, $(SRC_PLUGIN))


## @DIR_CONF_LINUX Local user configuration directory Linux
DIR_CONF_LINUX=${HOME}/.config/gosh
## @DIR_PLUGIN_LINUX Local user plugin directory Linux
DIR_PLUGIN_LINUX=${HOME}/.config/gosh/plugins

GOPATH := $(shell go env GOPATH)
GOTOOLCHAIN=local
APP=gosh
MAIN=cmd/gosh.go

.PHONY: build-all run plugin
all: clean build-all plugin

clean:
	rm -rf $(DIR_PKG)
	rm -rf $(DIR_BUILD)

plugin: build-plugin
	mv $(DIR_PKG_PLUGIN) $(DIR_BUILD_PLUGIN)

run: all
	GOSH_PLUGINS_DIR=$(DIR_BUILD_PLUGIN) $(DIR_BUILD)/linux/$(APP)
build-all: build-linux
build-linux:
	mkdir -p $(DIR_BUILD)/linux
	GOOS=linux go build -o $(DIR_BUILD)/linux/$(APP) $(MAIN)

build-plugin: $(OBJ_PLUGIN)
	echo compiling

%: $(DIR_SRC_PLUGIN)/%.go
	go build -buildmode=plugin -o $(DIR_PKG_PLUGIN)/$@.so $<

install:
	mkdir -p $(DIR_PLUGIN)
	cp $(DIR_PKG)/plugins/* $(DIR_PLUGIN_LINUX)/
	go install $(MAIN)

uninstall:
	rm -rf $(DIR_PLUGIN_LINUX)
	rm ${GOPATH}/bin/gosh

