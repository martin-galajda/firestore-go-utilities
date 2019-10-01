SHELL := bash

export PROJECT_ROOT := $(shell pwd)/
export TARGET_DIR = $(GOPATH)/src/github.com/martin-galajda/firestore-go-utilities/main
export TARGET_EXECUTABLE_PATH = ./bin/cli
export PATH := ./bin:$(PATH)

install: 
	@echo "GOPATH=$(GOPATH)"
	mkdir -p $(TARGET_DIR)
	cp -r ./src/* $(TARGET_DIR)
	go install github.com/martin-galajda/firestore-go-utilities/main

build-cli:
	@echo Building CLI...
	go build -o $(TARGET_EXECUTABLE_PATH) ./cmd/cli

run-get-images: build-cli
	cli -command=get-images

run-make-labelbox-labels: build-cli
	cli -command=make-labelbox-labels

run-labelbox-annotations-to-validation: build-cli
	cli -command=labelbox-annotations-to-validation-annotations -input_path=./out/export-2019-09-20T08_12_10.802Z.json -out_dir=./out/validation-annotations
