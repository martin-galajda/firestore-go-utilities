package main

import (
	"cloud.google.com/go/firestore"
	"fmt"
	"github.com/martin-galajda/firestore-go-utilities/internal/fileutils"
	"golang.org/x/net/context"
	"log"
	"flag"
	"path"
	"time"

	"github.com/martin-galajda/firestore-go-utilities/internal/googleapis"
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
		"./data/collected-dataset",
		"Path to the output directory.",
	)



	flag.Parse()
}


// defined and parsed in parse-flags.go
var pathToFirebaseConfigFile, firebaseProjectID, datasetName, pathToOutputDir *string
var ctx context.Context
var firestoreClient *firestore.Client

func init() {
	parseCLIFlags()
	ctx = context.Background()

	firestoreClient = mustInitGoogleAPIClients(ctx)
}

const (
	PAGE_SIZE = 10
	NUM_OF_PAGES = 1000
)

func mustInitGoogleAPIClients(ctx context.Context) (*firestore.Client) {
	firestoreClient, err := googleapis.NewFirestoreClient(ctx, *firebaseProjectID, *pathToFirebaseConfigFile)

	if err != nil {
		panic(err)
	}

	return firestoreClient
}

type firestoreUrlAnnotation struct {
	Metadata struct {
		AnnotatedElementsData map[string]struct{
			DataAnnotationId string
			ElemPathFromRoot string
			ImgUrlBase64 string
			Url string
		}
		Html string
	}
	Url string
}

func main() {

	pathToCol := fmt.Sprintf("workSessions/%v/processedUrlsData", *datasetName)

	var cur *firestore.DocumentSnapshot
	pageFetched := 0

	elems := []*firestoreUrlAnnotation{}
	reachedEnd := false
	count  := 0

	for pageFetched < NUM_OF_PAGES && !reachedEnd {
		var q firestore.Query
		if cur != nil {
			q = firestoreClient.
				Collection(pathToCol).
				OrderBy("url", firestore.Asc).
				Limit(PAGE_SIZE).
				StartAfter(cur)
		} else {
			q = firestoreClient.
				Collection(pathToCol).
				OrderBy("url", firestore.Asc).
				Limit(PAGE_SIZE)
		}

		docsIt := q.Documents(context.Background())

		docsSnapshots, err := docsIt.GetAll()

		if err != nil {
			log.Fatalf("Error getting document snapshot: %v\n", err)
		}

		if len(docsSnapshots) > 0 {
			cur = docsSnapshots[len(docsSnapshots) - 1]
		}

		for _, s := range docsSnapshots {
			var elem firestoreUrlAnnotation
			err = s.DataTo(&elem)

			if err != nil {
				log.Fatalf("Error parsing document snapshot: %v\n", err)
			}

			elems = append(elems, &elem)
		}

		count += len(docsSnapshots)

		reachedEnd = len(docsSnapshots) != PAGE_SIZE
		pageFetched++
	}


	outFilePath := path.Join(*pathToOutputDir, "export-" + time.Now().UTC().Format(time.RFC3339) + ".json")
	fileutils.CreateFile(&outFilePath)

	err := fileutils.WriteJSON(&outFilePath, elems)

	if err !=  nil {
		log.Fatalf("fileutils.WriteJSON: error writing exported JSON data: %v\n", err)
	}

	fmt.Printf("Successfully exported dataset from Firestore to '%v'\n", outFilePath)
	fmt.Printf("Exported %v processed urls\n", count)
}
