# Utility programs

## Purpose:
- moving data between Firestore and Labelbox
- import and export labels from Labelbox
- importing predictions to Labelbox via API
- exporting annotations from Labelbox


## Setup
1. Install go version 1.13 from https://golang.org
2. Generate service account key for Firestore and place it in `./.secrets/service-account.json`
3. Generate secret API token for Labelbox and place it in `./.secrets/labelbox-api-token.txt`
4. (optional) install `make` if you want to use commands from Makefile

## Included utilities:
- `cmd/cli` 
    - get data from Firestore
    - transform Open Images labels to format needed for Labelbox import
    - transform Open Images annotations to format needed by our evaluation scripts
- `cmd/compute-dataset-stats` - compute stats about collected dataset
- `cmd/construct-label-mapping` - debug reconstruction of labels
- `cmd/download-openimages-test-images` - download subset of test images from Open Images
- `cmd/export-dataset-datarows` - export rows from Labelbox
- `cmd/format-openiamges-test-images-annotations` - format GT boxes provided by Open Images into files with format needed for our evaluation script
- `cmd/init-prediction-model` -  initialize prediction model in Labelbox using API
- `cmd/make-prediction-model` -  import predictions made by models into Labelbox using API
- `cmd/switch-prediction-model` - switch active prediction model in Labelbox using API

## Running programs

Note: you need to have secrets in place to run these programs.

1. Run `make run-make-labelbox-labels` from the root to transform labels from Open Images to JSON file that can be imported to Labelbox when creating new project for annotating.
2. Run `make run-labelbox-annotations-to-validation` to transform manually exported annotations from Labelbox to format that is needed for our scripts.
3. Run `make run-get-images` to download annotated images from Firestore.
4. Run `go run cmd/export-firestore-dataset/main.go` to export collected dataset from Firestore to /data/collected-dataset.

Every program in `cmd` directory can be ran as `go run cmd/<name_of_the_program>` from the root directory.
