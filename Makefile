SHELL := bash

export GOPATH := $(shell pwd)/workspace
export TARGET_DIR = $(GOPATH)/src/github.com/martinGalajdaSchool/go-cli
export PATH := $(GOPATH)/bin:$(PATH)

install: 
	@echo "GOPATH=$(GOPATH)"
	mkdir -p $(TARGET_DIR)
	cp -r ./src/* $(TARGET_DIR)
	go install github.com/martinGalajdaSchool/go-cli

run-get-images:
	go-cli -command=get-images

run-make-labelbox-labels:
	go-cli -command=make-labelbox-labels
