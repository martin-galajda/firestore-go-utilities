package main

import (
	"flag"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/translate"
	"golang.org/x/net/context"
)

var pathToFirebaseConfigFile, firebaseProjectID, command, datasetName *string
var pathToOutputDir, pathToOutputFile, pathToInputsFile, pathToLabelboxLabelsOutputFile *string
var ctx context.Context
var firestoreClient *firestore.Client
var translateClient *translate.Client

func init() {
	initFlags()
	ctx, firestoreClient, translateClient = initGoogleApiClients()
}

func main() {
	switch *command {
	case "get-images":
		getImages(ctx, firestoreClient, *datasetName)
	case "make-labelbox-labels":
		transformLabelsToLabelBoxFormat(*pathToInputsFile, *pathToLabelboxLabelsOutputFile)
	case "labelbox-annotations-to-validation-annotations":
		transformLabelboxAnnotations(*pathToInputsFile, *pathToOutputDir)
	default:
		flag.PrintDefaults()
		log.Fatalf("Invalid command for CLI: %v.", command)
	}

	fmt.Println("Successfully executed CLI command!")
}
