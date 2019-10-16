# Utility programs

## Purpose:
- moving data between Firestore and Labelbox
- import and export labels from Labelbox
- importing predictions to Labelbox via API
- exporting annotations from Labelbox

Tiny Go programs:
- `cmd/cli` 
    - get data from Firestore
    - transform Open Images labels to format needed for Labelbox import
    - transform Open Images annotations to format needed by our evaluation scripts
- `cmd/compute-dataset-stats` - compute stats about collected dataset
- `cmd/construct-label-mapping` - debug reconstruction of labels
- `cmd/download-openimages-test-images` - download subset of test images from Open Images
- `cmd/export-dataset-datarows` - export rows from from Labelbox
- `cmd/format-openiamges-test-images-annotations` - format GT boxes provided by Open Images into files with format needed for our evaluation script
- `cmd/init-prediction-model` -  initialize prediction model in Labelbox using API
- `cmd/make-prediction-model` -  import predictions made by models into Labelbox using API
- `cmd/switch-prediction-model` - switch active prediction model in Labelbox using API
