package main

import (
	"flag"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"golang.org/x/net/context"

	"github.com/martin-galajda/firestore-go-utilities/internal/commands"
	"github.com/martin-galajda/firestore-go-utilities/internal/googleapis"
)

var pathToFirebaseConfigFile, firebaseProjectID, command, datasetName *string
var pathToOutputDir, pathToInputsFile, pathToLabelboxLabelsOutputFile *string
var ctx context.Context
var firestoreClient *firestore.Client
var translator googleapis.Translator

func init() {
	parseCLIFlags()
	ctx = context.Background()

	firestoreClient, translator = mustInitGoogleAPIClients(ctx)
}

func mustInitGoogleAPIClients(ctx context.Context) (*firestore.Client, googleapis.Translator) {

	firestoreClient, err := googleapis.NewFirestoreClient(ctx, *firebaseProjectID, *pathToFirebaseConfigFile)

	if err != nil {
		panic(err)
	}

	translator := googleapis.NewTranslator(ctx, *pathToFirebaseConfigFile)

	return firestoreClient, translator
}

func main() {
	switch *command {
	case "get-images":
		commands.GetImages(ctx, firestoreClient, *datasetName, *pathToOutputDir)
	case "make-labelbox-labels":
		commands.TransformLabelsToLabelBoxFormat(translator, *pathToInputsFile, *pathToLabelboxLabelsOutputFile)
	case "labelbox-annotations-to-validation-annotations":
		commands.TransformLabelboxAnnotations(*pathToInputsFile, *pathToOutputDir)
	default:
		flag.PrintDefaults()
		log.Fatalf("Invalid command for CLI: %q.", *command)
	}

	fmt.Println("Successfully executed CLI command!")
}
