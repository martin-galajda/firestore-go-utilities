SHELL := bash

export PROJECT_ROOT := $(shell pwd)/
export TARGET_DIR = $(GOPATH)/src/github.com/martin-galajda/firestore-go-utilities/main
export TARGET_EXECUTABLE_PATH = ./bin/
export PATH := ./bin:$(PATH)

install: 
	@echo "GOPATH=$(GOPATH)"
	mkdir -p $(TARGET_DIR)
	cp -r ./src/* $(TARGET_DIR)
	go install github.com/martin-galajda/firestore-go-utilities/main

build-cli:
	@echo Building CLI...
	go build -o $(TARGET_EXECUTABLE_PATH) ./cmd/cli

build-all:
	@echo Building commands...
	go build -o $(TARGET_EXECUTABLE_PATH) ./cmd/...

run-get-images: build-cli
	cli -command=get-images

run-make-labelbox-labels: build-cli
	cli -command=make-labelbox-labels

run-labelbox-annotations-to-validation: build-cli
	cli -command=labelbox-annotations-to-validation-annotations -input_path=./out/export-2019-10-10T16_08_11.790Z.json -out_dir=./out/validation-annotations

run-export-labelbox-dataset-rows: build-all
	export-dataset-datarows

