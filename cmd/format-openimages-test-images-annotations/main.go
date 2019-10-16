package main

import (
	"encoding/csv"
	"fmt"
	"flag"
	"image"
	_ "image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/martin-galajda/firestore-go-utilities/internal/annotation/labelutils"
)

var fPathToCSVFileWithImages *string
var fPathToCSVFileWithAnnotations *string
var fExportDir *string
var fImagesDir *string
var fImagesCount *int

const (
	csvIndexImageID   = 0
	csvIndexLabelName = 2
	csvIndexXMin      = 4
	csvIndexXMax      = 5
	csvIndexYMin      = 6
	csvIndexYMax      = 7
)

// file exported via Labelbox GUI
const pathToExportedLabelboxFile = "./data/class-descriptions-labelbox.json"

// file provided by OpenImages team
const pathToLabelsCSVFile = "./data/class-descriptions.csv"


// ImageID,Source,LabelName,Confidence,XMin,XMax,YMin,YMax,IsOccluded,IsTruncated,IsGroupOf,IsDepiction,IsInside
type annotationRecordFromCSV struct {
	ImageID   string  // 0
	LabelName string  // 2
	XMin      float32 // 4
	XMax      float32 //  5
	YMin      float32 // 6
	YMax      float32 // 7
}

func MustGetFloat(str string) float32 {
	v, err := strconv.ParseFloat(str, 32)

	if err != nil {
		panic(err)
	}

	return float32(v)
}

func loadCSVFileWithAnnotations(pathToCSVFile string, imagesCount int) map[string][]annotationRecordFromCSV {
	f, err := os.Open(pathToCSVFile)
	defer f.Close()

	log.Println("Opened CSV file")

	if err != nil {
		log.Fatalf("Error opening file for loading images: %s", err.Error())
	}

	csvReader := csv.NewReader(f)

	var parsedRecord annotationRecordFromCSV
	result := map[string][]annotationRecordFromCSV{}

	lineIdx := 0
	for {
		record, err := csvReader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Error reading CSV containing Open Images annotations: %s", err.Error())
		}

		if lineIdx == 0 {
			lineIdx++
			continue
		}

		parsedRecord.ImageID = record[csvIndexImageID]
		parsedRecord.LabelName = record[csvIndexLabelName]
		parsedRecord.XMin = MustGetFloat(record[csvIndexXMin])
		parsedRecord.XMax = MustGetFloat(record[csvIndexXMax])
		parsedRecord.YMin = MustGetFloat(record[csvIndexYMin])
		parsedRecord.YMax = MustGetFloat(record[csvIndexYMax])

		if _, ok := result[parsedRecord.ImageID]; !ok {
			result[parsedRecord.ImageID] = []annotationRecordFromCSV{}
		}

		result[parsedRecord.ImageID] = append(result[parsedRecord.ImageID], parsedRecord)
		lineIdx++
	}

	return result
}

func init() {
	parseFlags()
}

func parseFlags() {
	fPathToCSVFileWithAnnotations = flag.String(
		"path_to_annotations_csv",
		"./data/openimages-resources/test-images-annotations-bbox.csv",
		"Path to the CSV file containing annotations for test images",
	)

	fImagesDir = flag.String(
		"images_dir",
		"./data/openimages-resources/openimages-test-images/",
		"Path the export directory for images.",
	)

	fExportDir = flag.String(
		"export_dir",
		"./data/openimages-resources/ground_truth_openimages/",
		"Path the export directory for images.",
	)

	fImagesCount = flag.Int(
		"images_count",
		100,
		"No. of images to download.",
	)
}

func getImageDimensions(pathToImage string) (int, int, error) {
	reader, err := os.Open(pathToImage)

	log.Println(pathToImage)
	if err != nil {
		return 0, 0, err
	}
	defer reader.Close()
	im, _, err := image.DecodeConfig(reader)

	if err != nil {
		return 0, 0, err
	}
	return im.Width, im.Height, nil
}

func produceGTFile(pathToGTDirectory string, imageFilename string, midToLabel map[string]labelutils.AnnotationLabel, imageAnnotations []annotationRecordFromCSV, imWidth, imHeight int) error {
	filename := pathToGTDirectory + imageFilename + ".txt"

	f, err := os.Create(filename)

	if err != nil {
		return err
	}

	for _, annotation := range imageAnnotations {

		label := midToLabel[annotation.LabelName]

		minX := annotation.XMin * float32(imWidth)
		maxX := annotation.XMax * float32(imWidth)
		minY := annotation.YMin * float32(imHeight)
		maxY := annotation.YMax * float32(imHeight)
		fmt.Fprintf(f, "%s %d %d %d %d\n", label.EvaluationFormat, int(minX), int(minY), int(maxX), int(maxY))
	}

	return nil
}

func main() {
	log.Println(*fImagesCount)
	recordsFromCSV := loadCSVFileWithAnnotations(*fPathToCSVFileWithAnnotations, *fImagesCount)

	filesInDir, err := ioutil.ReadDir(*fImagesDir)

	if err != nil {
		panic(err)
	}

	imageExtensions := map[string]struct{}{
		".jpg":  struct{}{},
		".jpeg": struct{}{},
	}

	labels := labelutils.GetLabels(pathToLabelsCSVFile, pathToExportedLabelboxFile)

	midToLabel := map[string]labelutils.AnnotationLabel{}
	for _, label := range labels {
		midToLabel[label.Mid] = label
	}

	for _, fileInExportDir := range filesInDir {
		if _, ok := imageExtensions[path.Ext(fileInExportDir.Name())]; !ok {
			continue
		}

		imWidth, imHeight, err := getImageDimensions(filepath.Join(*fImagesDir, fileInExportDir.Name()))

		if err != nil {
			log.Printf("Error getting image dimensions: %s\n", err.Error())
			continue
		}
		imID := fileInExportDir.Name()
		imID = strings.Replace(imID, path.Ext(imID), "", 1)

		annotations := recordsFromCSV[imID]

		produceGTFile(*fExportDir, fileInExportDir.Name(), midToLabel, annotations, imWidth, imHeight)

		log.Printf("imID = %s\n", imID)
		log.Printf("imWidth=%d, imHeight=%d\n", imWidth, imHeight)
	}

	log.Printf("Successfully executed program: %d\n", len(recordsFromCSV))

}
