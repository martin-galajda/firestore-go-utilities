package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"golang.org/x/net/context"

	"github.com/martin-galajda/firestore-go-utilities/internal/fileutils"
	labelboxAPI "github.com/martin-galajda/firestore-go-utilities/internal/labelbox/api"
)

var fProjectID, fPathToLabelboxAPITokenFile, fPredictionModelID *string
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
		"prediction_model_id",
		// Faster R-CNN
		// "ck1k9et1ye1en083855rvodlz",

		// yolov3
		"ck1kvycemocj2079444c9o3h2",
		"Name of the prediction model.",
	)

	fPathToLabelboxAPITokenFile = flag.String(
		"path_to_labelbox_api_token",
		"./.secrets/labelbox-api-token.txt",
		"Path to the file containing secret API token for Labelbox.",
	)

	flag.Parse()
}

func mustGetLabelboxAPIToken(pathToTokenFile string) string {
	apiToken, err := fileutils.ReadFileAsText(pathToTokenFile)

	if err != nil {
		log.Fatalf("Error reading API token for LabelBox: %q", err.Error())
	}

	fmt.Printf("Successfully retrieved labelbox API token.\n")

	return strings.Replace(apiToken, "\n", "", 0)
}

func main() {
	ctx, _ = context.WithTimeout(ctx, time.Second*20)
	fmt.Print(labelboxAPIToken)

	projectID := *fProjectID
	predictionModelID := *fPredictionModelID

	labelboxAPIToken := strings.Replace(labelboxAPIToken, "\n", "", 1)
	apiClient := labelboxAPI.NewLabelboxAPI(labelboxAPIURL, labelboxAPIToken)

	err := apiClient.AttachPredictionModel(ctx, projectID, predictionModelID)

	if err != nil {
		log.Fatalf("Error creating prediction model: %s", err.Error())
	}
	fmt.Printf("New prediction model ID: %q\n", predictionModelID)

	fmt.Println("Successfully executed program!")
}
