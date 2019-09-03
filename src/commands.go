package main

import (
	"bufio"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
)

func getImages(ctx context.Context, client *firestore.Client, dataset string) {
	worksessionColRef := client.Collection("workSessions")
	datasetWorkSessionDocRef := worksessionColRef.Doc(dataset)
	processedURLColRef := datasetWorkSessionDocRef.Collection("processedUrls")

	documentRefIterator := processedURLColRef.DocumentRefs(ctx)

	if documentRefIterator == nil {
		log.Fatalf("Document ref iterator is nil for processedUrls. Dataset = %q", dataset)
	}

	urlsToDownload := []string{}

	count := 0
	maxCount := 1

	for {
		count++

		docRef, errRefIterator := documentRefIterator.Next()
		if errRefIterator == iterator.Done {
			break
		}

		if errRefIterator != nil {
			log.Printf("Error getting document ref. Err: %v", errRefIterator)
			continue
		}

		if count > maxCount {
			break
		}

		docSnapshot, errRef := docRef.Get(ctx)
		if errRef != nil {
			log.Printf("Error getting document snapshot. Err: %v", errRef)
			continue
		}

		var processedURLDoc ProcessedUrlDocument
		errSnapshot := docSnapshot.DataTo(&processedURLDoc)

		if errSnapshot != nil {
			log.Printf("Error getting processed url document data. Err: %v", errSnapshot)
			continue
		}

		for _, annotatedElementsData := range processedURLDoc.Data.AnnotatedElementsData {
			log.Printf("Collected new url: %q.\n", annotatedElementsData.Url)
			urlsToDownload = append(urlsToDownload, annotatedElementsData.Url)
		}
	}

	httpClient := makeHttpClient()

	for _, url := range urlsToDownload {
		downloadFile(httpClient, &url, *pathToOutputDir)
	}
}

func transformLabelsToLabelBoxFormat(pathToLabelsFile, pathToOutputFile string) {
	f, err := os.Open(pathToLabelsFile)

	if err != nil {
		log.Fatalf("Error opening file for transforming labels into labelbox format. Filepath: %q, error: %v", pathToLabelsFile, err)
	}

	fileReader := bufio.NewScanner(f)

	res := NewLabelboxLabelSettings()

	for fileReader.Scan() {
		line := fileReader.Text()

		// ignore empty lines
		if strings.TrimPrefix(line, " ") == "" {
			continue
		}

		tokens := strings.Split(line, ",")

		if len(tokens) != 2 {
			log.Fatalf("Expected lines in labels CSV file to contain one comma. Got %d. Tokens = %v", len(tokens), tokens)
		}

		mid, label := tokens[0], tokens[1]
		res.AddToolDefinition(mid, label)
	}

	err = writeJSON(&pathToOutputFile, res)

	if err != nil {
		log.Fatalf("Error writing labels for Labelbox to JSON file: %v", err)
	}
}
