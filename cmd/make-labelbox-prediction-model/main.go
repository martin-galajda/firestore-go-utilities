package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/martin-galajda/firestore-go-utilities/internal/annotation/labelutils"
	"github.com/martin-galajda/firestore-go-utilities/internal/fileutils"
	"github.com/martin-galajda/firestore-go-utilities/internal/labelbox"
	labelboxAPI "github.com/martin-galajda/firestore-go-utilities/internal/labelbox/api"
)

// file exported via Labelbox GUI
const pathToExportedLabelboxFile = "./data/class-descriptions-labelbox.json"

// file provided by OpenImages team
const pathToLabelsCSVFile = "./data/class-descriptions.csv"

var fProjectID, fPathToLabelboxAPITokenFile, fPredictionModelID, fDatasetID, fModelPredictions *string
var ctx context.Context
var labelboxAPIToken string

const labelboxAPIURL = "https://api.labelbox.com/graphql"

func init() {
	parseCLIFlags()
	ctx = context.Background()

	labelboxAPIToken = mustGetLabelboxAPIToken(*fPathToLabelboxAPITokenFile)
}

func parseCLIFlags() {
	fProjectID = flag.String(
		"project_id",
		"ck1k95vi15tds0846qjcy5nli",
		"ID of the project in Labelbox.",
	)

	fPredictionModelID = flag.String(
		"prediction_model_ID",
		//"ck1k9et1ye1en083855rvodlz",
		"ck1kvycemocj2079444c9o3h2",
		"ID of the prediction model in Labelbox.",
	)

	fDatasetID = flag.String(
		"dataset_id",
		"ck0c3qn7i5gp80863xprnsnu4",
		"ID of the dataset in Labelbox.",
	)

	fPathToLabelboxAPITokenFile = flag.String(
		"path_to_labelbox_api_token",
		"./.secrets/labelbox-api-token.txt",
		"Path to the file containing secret API token for Labelbox.",
	)

	fModelPredictions = flag.String(
		"model_predictions",
		"FRCNN",
		"Type of the model predictions - one of 'FRCNN', 'YOLOv3'.",
	)


	flag.Parse()
}

func getPathToModelDetections(modelForPredictions string) string {
	if modelForPredictions == "FRCNN" {
		return "./out/detected/test-images/FasterRCNN"
	} else if modelForPredictions == "YOLOv3" {
		return "./out/detected/test-images/YOLOv3"
	}

	panic("Invalid model specified for predictions. Use one of 'YOLOv3', 'FRCNN'.")
}

func mustGetLabelboxAPIToken(pathToTokenFile string) string {
	apiToken, err := fileutils.ReadFileAsText(pathToTokenFile)

	if err != nil {
		log.Fatalf("Error reading API token for LabelBox: %q", err.Error())
	}

	fmt.Printf("Successfully retrieved labelbox API token.\n")

	return strings.Replace(apiToken, "\n", "", 0)
}

func loadModelDetections(pathToDirWithDetections string) map[string]map[string][]labelbox.LabelboxExportLabel {

	dirFilePaths, err := fileutils.GetFilePathsInDirWithExt(pathToDirWithDetections, ".txt")

	if err != nil {
		log.Fatalf("Error reading directory containing model detections. Path: %s. Error: %s\n", pathToDirWithDetections, err.Error())
	}

	res := map[string]map[string][]labelbox.LabelboxExportLabel{}

	allLabelsList := labelutils.GetLabels(pathToLabelsCSVFile, pathToExportedLabelboxFile)
	labelByEvaluationFormat := map[string]labelutils.AnnotationLabel{}

	for _, labelVal := range allLabelsList {
		labelByEvaluationFormat[labelVal.EvaluationFormat] = labelVal
	}

	fmt.Println(dirFilePaths)

	for _, pathToFileWithDetections := range dirFilePaths {
		f, err := os.Open(pathToFileWithDetections)

		if err != nil {
			log.Fatalf("Error reading file containing model detections. Path: %s. Error: %s\n", pathToFileWithDetections, err.Error())
		}


		resForFile := map[string][]labelbox.LabelboxExportLabel{}


		labelExternalID := strings.Replace(path.Base(pathToFileWithDetections), ".txt", "", 1)
		res[labelExternalID] = resForFile
		buffReader := bufio.NewScanner(f)

		for buffReader.Scan() {
			line := buffReader.Text()
			fmt.Println(line)

			if strings.Trim(line, " ") == "" {
				continue
			}

			tokens := strings.Split(line, " ")

			asInt32 := func(s string) int32 {
				val, err := strconv.Atoi(s)

				if err != nil {
					log.Fatalf("Invalid format for detected box coordinate: %v\n", err.Error())
				}

				return int32(val)
			}

			evaluationLabel := tokens[0]
			topLeftX, topLeftY := asInt32(tokens[2]), asInt32(tokens[3])
			bottomRightX, bottomRightY := asInt32(tokens[4]), asInt32(tokens[5])

			topLeftPoint := labelbox.LabelboxExportLabelGeometry{
				X: topLeftX,
				Y: topLeftY,
			}

			bottomRightPoint := labelbox.LabelboxExportLabelGeometry{
				X: bottomRightX,
				Y: bottomRightY,
			}

			topRightPoint := labelbox.LabelboxExportLabelGeometry{
				X: bottomRightX,
				Y: topLeftPoint.Y,
			}

			bottomLeftPoint := labelbox.LabelboxExportLabelGeometry{
				X: topLeftPoint.X,
				Y: bottomRightY,
			}

			labelboxLabel := labelbox.LabelboxExportLabel{
				Geometry: []labelbox.LabelboxExportLabelGeometry{topLeftPoint, topRightPoint, bottomRightPoint, bottomLeftPoint},
			}

			labelID := labelByEvaluationFormat[evaluationLabel].LabelboxHumanReadable
			if _, ok := resForFile[labelID]; !ok {
				resForFile[labelID] = []labelbox.LabelboxExportLabel{}
			}

			resForFile[labelID] = append(resForFile[labelID], labelboxLabel)
		}
	}


	fmt.Printf("len(res)=%d\n", len(res))

	return res
}

func main() {

	pathToDetections := getPathToModelDetections(*fModelPredictions)
	detections := loadModelDetections(pathToDetections)

	apiClient := labelboxAPI.NewLabelboxAPI(labelboxAPIURL, strings.Replace(labelboxAPIToken, "\n", "", 1))
	rows, err := apiClient.ExportProjectDatasetRows(ctx, *fProjectID, []string{*fDatasetID})

	fmt.Printf("Retrieved %d dataset rows.", len(rows))

	if err != nil {
		log.Fatalf("Error getting dataset rows: %s", err.Error())
	}

	datasetRowByExternalId := map[string]labelboxAPI.GraphQLDatasetRow{}

	for _, row := range rows {
		datasetRowByExternalId[row.ExternalID] = row
	}

	projectID := *fProjectID
	predictionModelID := *fPredictionModelID
	for externalLabelID, detection := range detections {
		datasetRow := datasetRowByExternalId[externalLabelID]

		err = apiClient.CreatePrediction(ctx, detection, projectID, predictionModelID, datasetRow.ID)

		if err != nil {
			log.Fatalf("Error creating prediction: %s", err.Error())
		}
	}

	fmt.Println("Sucessfully imported predictions!")
}
