package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

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
	fileIds := []string{}
	filenamesMetadata := [][]string{}

	for {
		docRef, errRefIterator := documentRefIterator.Next()
		if errRefIterator == iterator.Done {
			break
		}

		if errRefIterator != nil {
			log.Printf("Error getting document ref. Err: %v", errRefIterator)
			continue
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
			log.Printf("Collected new url: %q. Data  annotation ID: %s.\n", annotatedElementsData.Url, annotatedElementsData.DataAnnotationID)

			fileID := annotatedElementsData.DataAnnotationID
			if fileID == "" {
				fileID = makeUUID()
			}

			urlsToDownload = append(urlsToDownload, annotatedElementsData.Url)
			fileIds = append(fileIds, fileID)
			filenamesMetadata = append(filenamesMetadata, []string{fileID, annotatedElementsData.Url})
		}
	}

	httpClient := makeHttpClient()

	destinationDirForFiles := *pathToOutputDir + "/export-" + time.Now().Format(time.RFC3339)
	createDirIfNotExists(destinationDirForFiles)
	for idx, url := range urlsToDownload {
		downloadAndSaveFile(httpClient, &url, fileIds[idx], destinationDirForFiles)
	}

	metadataFilenamePath := path.Join(destinationDirForFiles, "filenames_metadata.csv")
	metadataFilenamesFile := createFile(&metadataFilenamePath)

	csv.NewWriter(metadataFilenamesFile).WriteAll(filenamesMetadata)
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

func transformLabelboxAnnotations(pathToLabelboxAnnotationsFile, pathToOutputDir string) {
	f, err := os.Open(pathToLabelboxAnnotationsFile)

	if err != nil {
		log.Fatalf("Error opening file for transforming exported labelbox annotations. Filepath: %q, error: %v", pathToLabelboxAnnotationsFile, err)
	}

	fileBytes, err := ioutil.ReadAll(f)

	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	exportedAnnotations := []LabelboxExportAnnotation{}
	json.Unmarshal(fileBytes, &exportedAnnotations)

	rowsForFiles := map[string][][]string{}

	// regex that will replace translated part of the label - the stuff in brackets e.g -> Human(Clovek) --> Human
	labelTranslationRe := regexp.MustCompile(`\(.+\)$`)
	for _, exportedAnnotation := range exportedAnnotations {
		fileID := exportedAnnotation.ExternalID
		rows := [][]string{}

		for class, classLabels := range exportedAnnotation.Labels {
			classWithoutTranslation := labelTranslationRe.ReplaceAllString(class, "")
			classWithoutTranslation = regexp.MustCompile(`\s`).ReplaceAllString(classWithoutTranslation, "")
			classWithoutTranslation = strings.ToLower(classWithoutTranslation)

			for _, labelGeometry := range classLabels {
				fmt.Println(classWithoutTranslation)

				topLeft := labelGeometry.Geometry[0]
				bottomRight := labelGeometry.Geometry[3]
				str := fmt.Sprint
				row := []string{classWithoutTranslation, str(topLeft.X), str(topLeft.Y), str(bottomRight.X), str(bottomRight.Y)}
				rows = append(rows, row)
			}
		}

		rowsForFiles[fileID] = rows
	}

	err = createDirIfNotExists(pathToOutputDir)
	if err != nil {
		log.Panicf("Failed to create output directory for transformed labelbox annotations. Error: %v", err)
	}

	for fileID, rows := range rowsForFiles {
		fmt.Println(fileID)
		fmt.Println(rows)
		outFilePath := path.Join(pathToOutputDir, fileID+".txt")
		outFile := createFile(&outFilePath)
		csvWriter := csv.NewWriter(outFile)
		csvWriter.Comma = rune(" "[0])
		csvWriter.WriteAll(rows)
	}
}
