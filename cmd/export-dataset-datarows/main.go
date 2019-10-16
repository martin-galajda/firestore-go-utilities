package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"golang.org/x/net/context"

	"github.com/martin-galajda/firestore-go-utilities/internal/fileutils"
	labelboxAPI "github.com/martin-galajda/firestore-go-utilities/internal/labelbox/api"
)

var projectID, datasetID *string
var pathToOutputFile *string
var ctx context.Context
var pathToLabelboxAPITokenFile *string
var labelboxAPIToken string

const labelboxAPIURL = "https://api.labelbox.com/graphql"

func init() {
	parseCLIFlags()
	ctx = context.Background()

	labelboxAPIToken = mustGetLabelboxAPIToken(*pathToLabelboxAPITokenFile)
}

func parseCLIFlags() {
	pathToOutputFile = flag.String(
		"out_file",
		"./out/dataset[id=%s]-datarows",
		"Path to the output file.",
	)

	projectID = flag.String(
		"project_id",
		"ck05pb93x1igo0848opn6889e",
		"ID of the project in Labelbox.",
	)

	datasetID = flag.String(
		"dataset_id",
		"ck0c3qn7i5gp80863xprnsnu4",
		"ID of the dataset in Labelbox.",
	)

	pathToLabelboxAPITokenFile = flag.String(
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
	ctx, _ = context.WithTimeout(ctx, time.Second*10)
	fmt.Print(labelboxAPIToken)

	apiClient := labelboxAPI.NewLabelboxAPI(labelboxAPIURL, strings.Replace(labelboxAPIToken, "\n", "", 1))
	rows, err := apiClient.ExportProjectDatasetRows(ctx, *projectID, []string{*datasetID})

	if err != nil {
		log.Fatalf("Error getting dataset row from Labelbox: %q\n", err.Error())
	}

	rowsForOutput := [][]string{[]string{"labelbox_id", "external_id"}}
	for _, row := range rows {
		rowData := []string{row.ID, row.ExternalID}
		rowsForOutput = append(rowsForOutput, rowData)
	}

	*pathToOutputFile = fmt.Sprintf(*pathToOutputFile, *datasetID)
	outFile := fileutils.CreateFile(pathToOutputFile)
	csvWriter := csv.NewWriter(outFile)
	csvWriter.WriteAll(rowsForOutput)

	fmt.Println("Successfully executed program!")
}
