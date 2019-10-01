package commands

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/martin-galajda/firestore-go-utilities/internal/fileutils"
	"github.com/martin-galajda/firestore-go-utilities/internal/httputils"
	"github.com/martin-galajda/firestore-go-utilities/internal/labelbox"
	"github.com/martin-galajda/firestore-go-utilities/internal/uuid"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	googleFirestore "cloud.google.com/go/firestore"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
	"github.com/martin-galajda/firestore-go-utilities/internal/googleapis"
	"github.com/martin-galajda/firestore-go-utilities/internal/firestore"
)

// GetImages gets images marked as important from Firestore database.
func GetImages(ctx context.Context, client *googleFirestore.Client, dataset, pathToOutputDir string) {
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

		var processedURLDoc firestore.ProcessedUrlDocument
		errSnapshot := docSnapshot.DataTo(&processedURLDoc)

		if errSnapshot != nil {
			log.Printf("Error getting processed url document data. Err: %v", errSnapshot)
			continue
		}

		for _, annotatedElementsData := range processedURLDoc.Data.AnnotatedElementsData {
			log.Printf("Collected new url: %q. Data  annotation ID: %s.\n", annotatedElementsData.Url, annotatedElementsData.DataAnnotationID)

			fileID := annotatedElementsData.DataAnnotationID
			if fileID == "" {
				fileID = uuid.MakeUUID()
			}

			urlsToDownload = append(urlsToDownload, annotatedElementsData.Url)
			fileIds = append(fileIds, fileID)
			filenamesMetadata = append(filenamesMetadata, []string{fileID, annotatedElementsData.Url})
		}
	}

	httpClient := httputils.NewHTTPClient()

	destinationDirForFiles := pathToOutputDir + "/export-" + time.Now().Format(time.RFC3339)
	fileutils.CreateDirIfNotExists(destinationDirForFiles)
	for idx, url := range urlsToDownload {
		fileutils.DownloadAndSaveFile(httpClient, &url, fileIds[idx], destinationDirForFiles)
	}

	metadataFilenamePath := path.Join(destinationDirForFiles, "filenames_metadata.csv")
	metadataFilenamesFile := fileutils.CreateFile(&metadataFilenamePath)

	csv.NewWriter(metadataFilenamesFile).WriteAll(filenamesMetadata)
}

// TransformLabelsToLabelBoxFormat transforms OpenImages labels(=classes)
// to format needed for importing labels into LabelBox.
func TransformLabelsToLabelBoxFormat(translator googleapis.Translator, pathToLabelsFile, pathToOutputFile string) {
	f, err := os.Open(pathToLabelsFile)

	if err != nil {
		log.Fatalf("Error opening file for transforming labels into labelbox format. Filepath: %q, error: %v", pathToLabelsFile, err)
	}

	fileReader := bufio.NewScanner(f)

	res := labelbox.NewLabelboxLabelSettings()

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

		translatedLabel, err := translator.Translate(label)

		if err != nil {
			log.Fatalf("Error occurred when translating label: %v.", err)
		}

		res.AddToolDefinition(mid, label, translatedLabel)
	}

	err = fileutils.WriteJSON(&pathToOutputFile, res)

	if err != nil {
		log.Fatalf("Error writing labels for Labelbox to JSON file: %v", err)
	}
}

// TransformLabelboxAnnotations transforms exported annotations from Labelbox
// to format
func TransformLabelboxAnnotations(pathToLabelboxAnnotationsFile, pathToOutputDir string) {
	f, err := os.Open(pathToLabelboxAnnotationsFile)

	if err != nil {
		log.Fatalf("Error opening file for transforming exported labelbox annotations. Filepath: %q, error: %v", pathToLabelboxAnnotationsFile, err)
	}

	fileBytes, err := ioutil.ReadAll(f)

	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	exportedAnnotations := []labelbox.LabelboxExportAnnotation{}
	err = json.Unmarshal(fileBytes, &exportedAnnotations)

	if err != nil {
		log.Fatalf("Error Unmarshaling expoted annotations: %s", err)
	}

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

				leftTopPoint, rightBottomPoint := labelGeometry.GetBoundingBoxPoints()

				str := fmt.Sprint
				row := []string{
					classWithoutTranslation,
					str(leftTopPoint.X),
					str(leftTopPoint.Y),
					str(rightBottomPoint.X),
					str(rightBottomPoint.Y),
				}
				rows = append(rows, row)
			}
		}

		rowsForFiles[fileID] = rows
	}

	err = fileutils.CreateDirIfNotExists(pathToOutputDir)
	if err != nil {
		log.Panicf("Failed to create output directory for transformed labelbox annotations. Error: %v", err)
	}

	currentDirOutput := path.Join(pathToOutputDir, fileutils.GetFilenameFromFilepath(pathToLabelboxAnnotationsFile))
	err = fileutils.CreateDirIfNotExists(currentDirOutput)
	if err != nil {
		log.Panicf("Failed to create output directory for transformed labelbox annotations. Error: %v", err)
	}

	for fileID, rows := range rowsForFiles {
		fmt.Println(fileID)
		fmt.Println(rows)
		outFilePath := path.Join(currentDirOutput, fileID+".txt")
		outFile := fileutils.CreateFile(&outFilePath)
		csvWriter := csv.NewWriter(outFile)
		csvWriter.Comma = rune(" "[0])
		csvWriter.WriteAll(rows)
	}
}
