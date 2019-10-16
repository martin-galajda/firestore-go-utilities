package main

import (
	"github.com/martin-galajda/firestore-go-utilities/internal/fileutils"
	"fmt"
	"bufio"
	"strings"
	"log"
	"io/ioutil"
	"flag"
	"os"
)

var fPathToDirectoryWithGTBoxes *string
var fPathToDirectoryWithImages *string

func init() {
	parseFlags()
}


func parseFlags() {
	fPathToDirectoryWithGTBoxes = flag.String(
		"path_to_gt_directory",
		"../object-detection/evaluation/ground_truth/test-images",
		"Path to the directory containing annotated ground truth boxes.",
	)

	fPathToDirectoryWithImages = flag.String(
		"path_to_images_dir",
		"/Users/martingalajda/School/DIPLOMA-THESIS/firestore-go-utilities/out/images/export-2019-09-05T13:59:12+02:00",
		"Path to the directory containing images collected for annotation.",
	)
}

func readFileLines(pathToFile string) ([]string, error) {
	file, err := os.Open(pathToFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Trim(line, " ") != "" {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return lines, nil
}

func printDatasetStats(annotatedClassCount map[string]int, totalClassCount int, numOfCollectedImages, numOfAnnotatedImages int) {
	fmt.Println("*** Dataset Stats ***")

	fmt.Printf("Total bounding boxes: %d\n", totalClassCount)
	fmt.Printf("Average no. of bounding boxes: %0.2f\n", float32(totalClassCount)/float32(numOfAnnotatedImages))
	fmt.Printf("No. of collected important images: %d\n", numOfCollectedImages)
	fmt.Printf("No. of annotated important images: %d\n", numOfAnnotatedImages)
	fmt.Printf("No. of used classes: %d\n", len(annotatedClassCount))

	percentageCollectedByClass := map[string]float32{}

	otherCountNumber := totalClassCount
	for class, collectedCount := range annotatedClassCount  {
		percentageCollectedByClass[class] = float32(collectedCount) / float32(totalClassCount)
	}

	remaining := float32(100.)
	otherCount := len(annotatedClassCount)
	for class, collectedCountPercentage := range percentageCollectedByClass {
		if collectedCountPercentage > 0.025 {
			fmt.Printf("--- %s - %d - %.02f\n", class, annotatedClassCount[class], collectedCountPercentage*100)
			remaining -= collectedCountPercentage * 100
			otherCount -= 1
			otherCountNumber -= annotatedClassCount[class]
		}
	}

	fmt.Printf("--- Other(%d) - %d - %.02f\n", otherCount, otherCountNumber, remaining)
}

func main() {
	pathToGTDir := *fPathToDirectoryWithGTBoxes

	files, err := ioutil.ReadDir(*fPathToDirectoryWithImages)

	if err != nil {
		log.Fatalf("Error reading directory containing collected images: %s", err.Error())
	}

	numOfCollectedImages := len(files)

	gtAnnotationFilePaths, err := fileutils.GetFilePathsInDirWithExt(pathToGTDir, ".txt")

	if err != nil {
		log.Fatalf("Error reading directory containing ground truth annotations: %s", err.Error())
	}


	annotatedClassCount := map[string]int{}
	totalClassCount := 0

	for _, pathToGTFile := range gtAnnotationFilePaths {
		lines, err := readFileLines(pathToGTFile)

		if err != nil {
			log.Fatalf("Error reading file containing ground truth annotations: %s", err.Error())
		}

		for _, line := range lines {
			tokens := strings.Split(line, " ")
			class := tokens[0]

			if _, ok := annotatedClassCount[class]; !ok {
				annotatedClassCount[class] = 0
			}

			annotatedClassCount[class] += 1
			totalClassCount += 1
		}
	}

	printDatasetStats(
		annotatedClassCount,
		totalClassCount,
		numOfCollectedImages,
		len(gtAnnotationFilePaths),
	)


}
