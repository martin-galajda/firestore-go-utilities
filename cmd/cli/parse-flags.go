package main

import (
	"flag"
)

func parseCLIFlags() {
	pathToFirebaseConfigFile = flag.String(
		"firebase_config_file",
		"./.secrets/service-account.json",
		"Path to config file containing firebase service account credentials.",
	)

	firebaseProjectID = flag.String(
		"firebase_project_id",
		"download-images-plugin",
		"Name of the firebase project.",
	)

	datasetName = flag.String(
		"dataset_name",
		"output-2019-07-10T09:43:16-export-top-5000",
		"Name of the dataset to work with.",
	)

	pathToOutputDir = flag.String(
		"out_dir",
		"./out/images",
		"Path to the output directory.",
	)

	command = flag.String(
		"command",
		"",
		"Command to execute. One of the [get-images, make-labelbox-labels, labelbox-annotations-to-validation-annotations, export-labelbox-dataset-rows].",
	)

	pathToInputsFile = flag.String(
		"input_path",
		"../data/class-descriptions-boxable.csv",
		"Path to the file containing labels for OpenImages in form of original CSV file.",
	)

	pathToLabelboxLabelsOutputFile = flag.String(
		"out_path_to_labelbox_labels_file",
		"./out/class-descriptions-labelbox.json",
		"Path to the file containing labels for OpenImages in form of original CSV file.",
	)


	flag.Parse()
}
