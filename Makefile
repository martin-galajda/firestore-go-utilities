SHELL := bash

export GOPATH := $(shell pwd)/workspace
export TARGET_DIR = $(GOPATH)/src/github.com/martinGalajdaSchool/go-cli
export PATH := $(GOPATH)/bin:$(PATH)

install: 
	@echo "GOPATH=$(GOPATH)"
	mkdir -p $(TARGET_DIR)
	cp -r ./src/* $(TARGET_DIR)
	go install github.com/martinGalajdaSchool/go-cli

run-get-images: install
	go-cli -command=get-images

run-make-labelbox-labels: install
	go-cli -command=make-labelbox-labels

run-labelbox-annotations-to-validation: install
	go-cli -command=labelbox-annotations-to-validation-annotations -input_path=./out/export-2019-09-04T10_02_45.464Z.json -out_dir=./out/validation-annotations
