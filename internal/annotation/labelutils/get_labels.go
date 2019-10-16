package labelutils

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
)

const pathToLabelsFile = "./out/class-descriptions-labelbox.json"

type AnnotationLabel struct {
	Mid                   string
	OriginalHumanReadable string
	LabelboxHumanReadable string
	EvaluationFormat      string
	ExternalDatarowID     string
}

type LabelboxFormatTool struct {
	Mid  string
	Name string
}

type parsedDataFromLabelboxFormat struct {
	Tools []LabelboxFormatTool
	ExternalDataRowID string
}

type dataFromOriginalCSVFile struct {
	Mid                string
	HumanReadableClass string
	ExternalDataRowID  string
}

func getLabelsFromLabelboxFormatFile(pathToLabelsFile string) parsedDataFromLabelboxFormat {
	content, err := ioutil.ReadFile(pathToLabelsFile)

	if err != nil {
		log.Fatalf("Error reading file: %q", err.Error())
	}

	var fileData parsedDataFromLabelboxFormat

	err = json.Unmarshal(content, &fileData)

	if err != nil {
		log.Fatalf("Error unmarshaling data: %q", err.Error())
	}

	fileData.ExternalDataRowID = strings.Replace(path.Base(pathToLabelsFile), ".txt", "", 1)

	return fileData
}

func getLabelsFromOriginalCSVFile(pathToOriginalCSVFile string) []dataFromOriginalCSVFile {
	result := []dataFromOriginalCSVFile{}

	f, err := os.Open(pathToOriginalCSVFile)

	if err != nil {
		log.Fatalf("Error opening file for transforming labels into labelbox format. Filepath: %q, error: %v", pathToOriginalCSVFile, err)
	}

	fileReader := bufio.NewScanner(f)

	externalDatarowID := strings.Replace(path.Base(pathToOriginalCSVFile), ".txt", "", 1)

	for fileReader.Scan() {
		line := fileReader.Text()

		// ignore empty lines
		if strings.TrimPrefix(line, " ") == "" {
			continue
		}

		tokens := strings.Split(line, ",")

		mid := tokens[0]
		labelFromCSV := strings.Join(tokens[1:], ".")
		result = append(result, dataFromOriginalCSVFile{Mid: mid, HumanReadableClass: labelFromCSV, ExternalDataRowID: externalDatarowID })
	}

	return result
}

func mergeLabels(labelboxFormatByMid []LabelboxFormatTool, origFormatByMid []dataFromOriginalCSVFile) []AnnotationLabel {

	mergedByMid := map[string]*AnnotationLabel{}

	getEvaluationFormatVal := func(s string) string {
		s = regexp.MustCompile(`\(.+\)$`).ReplaceAllString(s, "")
		s = regexp.MustCompile(`\s`).ReplaceAllString(s, "")
		s = strings.ToLower(s)
		return s
	}

	for _, origFormatLabel := range origFormatByMid {
		mergedByMid[origFormatLabel.Mid] = &AnnotationLabel{
			Mid:                   origFormatLabel.Mid,
			OriginalHumanReadable: origFormatLabel.HumanReadableClass,
		}
	}

	for _, lbFormat := range labelboxFormatByMid {
		mergedByMid[lbFormat.Mid].LabelboxHumanReadable = lbFormat.Name
		mergedByMid[lbFormat.Mid].EvaluationFormat = getEvaluationFormatVal(lbFormat.Name)
	}

	resultValues := make([]AnnotationLabel, 0, len(mergedByMid))

	for _, val := range mergedByMid {
		resultValues = append(resultValues, *val)
	}

	return resultValues
}

// GetLabels gathers all information about labels from original CSV file
// and file containing exported annotations from Labelbox
// and merges them into one unified data structure with all the information
// such as original label name, transformed label name in Labelbox, format used for evaluation
// and its machine id (mid)
func GetLabels(pathToOriginalCSV, pathToLabelboxFormatFile string) []AnnotationLabel {
	dataFromLabelboxFormat := getLabelsFromLabelboxFormatFile(pathToLabelboxFormatFile)
	dataFromOriginalCSVFile := getLabelsFromOriginalCSVFile(pathToOriginalCSV)

	mergedLabels := mergeLabels(dataFromLabelboxFormat.Tools, dataFromOriginalCSVFile)

	return mergedLabels
}
