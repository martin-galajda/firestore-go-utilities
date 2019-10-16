package main

import (
	"encoding/csv"
	"image"
	_ "image/jpeg"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/martin-galajda/firestore-go-utilities/internal/httputils"
)

var fPathToCSVFileWithImages *string
var fPathToCSVFileWithAnnotations *string
var fExportDir *string
var fImagesCount *int

const (
	csvIndexImageID          = 0
	csvIndexOriginalURL      = 2
	csvIndexThumbnail300KURL = 10
)

type imageRecordFromCSV struct {
	ImageID          string // 0
	OriginalURL      string // 2
	Thumbnail300KURL string // 10
}

func loadCSVFileWithImages(pathToCSVFile string, imagesCount int) []imageRecordFromCSV {
	f, err := os.Open(pathToCSVFile)
	defer f.Close()

	log.Println("Opened CSV file")

	if err != nil {
		log.Fatalf("Error opening file for loading images: %s", err.Error())
	}

	csvReader := csv.NewReader(f)

	var parsedRecord imageRecordFromCSV
	result := []imageRecordFromCSV{}

	lineIdx := 0
	for lineIdx <= imagesCount {
		record, err := csvReader.Read()
		log.Println(record)

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Error reading CSV containing Open Images images: %s", err.Error())
		}

		if lineIdx == 0 {
			lineIdx++
			continue
		}

		parsedRecord.ImageID = record[csvIndexImageID]
		parsedRecord.OriginalURL = record[csvIndexOriginalURL]
		parsedRecord.Thumbnail300KURL = record[csvIndexThumbnail300KURL]

		result = append(result, parsedRecord)
		lineIdx++
	}

	return result
}

func init() {
	parseFlags()
}

func parseFlags() {
	fPathToCSVFileWithImages = flag.String(
		"path_to_images_csv",
		"./data/openimages-resources/test-images.csv",
		"Path to the CSV file containing metadata about test images",
	)

	fPathToCSVFileWithAnnotations = flag.String(
		"path_to_annotations_csv",
		"./data/openimages-resources/test-images-annotations-bbox.csv",
		"Path to the CSV file containing annotations for test images",
	)

	fExportDir = flag.String(
		"export_dir",
		"./data/openimages-resources/openimages-test-images/",
		"Path the export directory for images.",
	)

	fImagesCount = flag.Int(
		"images_count",
		1000,
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


func downloadAndSaveImage(client *http.Client, record imageRecordFromCSV, pathToOutputDir string) error {
	resp, err := client.Get(record.OriginalURL)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	ext := path.Ext(record.OriginalURL)
	filename := pathToOutputDir + record.ImageID + ext
	outFile, err := os.Create(filename)

	if err != nil {
		return err
	}

	defer outFile.Close()

	size, err := io.Copy(outFile, resp.Body)

	if err != nil {
		return err
	}

	fmt.Printf("Just Downloaded a file %s with size %d\n", filename, size)

	_, _, err = getImageDimensions(filename)

	if err != nil {
		os.Remove(filename)
		return err
	}

	return nil
}

func main() {
	log.Println(*fImagesCount)
	recordsFromCSV := loadCSVFileWithImages(*fPathToCSVFileWithImages, *fImagesCount)

	httpClient := httputils.NewHTTPClient()
	for _, recordFromCSV := range recordsFromCSV {
		err := downloadAndSaveImage(httpClient, recordFromCSV, *fExportDir)

		if err != nil {
			log.Printf("Error downloading image: %s", err.Error())
		}
	}

}
