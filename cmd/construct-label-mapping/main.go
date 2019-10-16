package main

import (
	"fmt"

	"github.com/martin-galajda/firestore-go-utilities/internal/annotation/labelutils"
)

// file exported via Labelbox GUI
const pathToExportedLabelboxFile = "./data/class-descriptions-labelbox.json"

// file provided by OpenImages team
const pathToLabelsCSVFile = "./data/class-descriptions.csv"

func main() {

	labels := labelutils.GetLabels(pathToLabelsCSVFile, pathToExportedLabelboxFile)

	for _, valLabel := range labels {
		fmt.Printf("valLabel.Mid=%v\n", valLabel.Mid)
		fmt.Printf("valLabel.OriginalHumanReadable=%v\n", valLabel.OriginalHumanReadable)
		fmt.Printf("valLabel.LabelboxHumanReadable=%v\n", valLabel.LabelboxHumanReadable)
		fmt.Printf("valLabel.EvaluationFormat=%v\n", valLabel.EvaluationFormat)
	}
}
