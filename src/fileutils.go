package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"

	"encoding/json"
	"path/filepath"
)

func downloadAndSaveFile(client *http.Client, url *string, fileId string, pathToOutputDir string) {
	resp, err := client.Get(*url)
	if err != nil {
		log.Fatalf("Error occurred when downloading image file from URL: %q. Error: %v", *url, err)
	}
	defer resp.Body.Close()

	filename, err := buildFilename(url, fileId)

	if err != nil {
		log.Printf("Failed to export file with url: %q\n. Skipping...", *url)
		return
	}

	fullPath := path.Join(pathToOutputDir, filename)

	file := createFile(&fullPath)

	size, err := io.Copy(file, resp.Body)
	defer file.Close()

	log.Printf("Downloaded and saved a file %q with size %d\n", filename, size)
}

func buildFilename(fullURL *string, fileID string) (string, error) {
	ext := filepath.Ext(*fullURL)

	if ext == "" {
		errMsg := fmt.Sprintf("Couldn't extract extension from URL: %q. Skipping import...", *fullURL)
		err := fmt.Errorf(errMsg)
		log.Printf("Couldn't extract extension from URL: %q. Skipping import...", *fullURL)

		return "", err
	}

	// get rid of query string from URL
	ext = regexp.MustCompile(`\?.+$`).ReplaceAllString(ext, "")

	return fileID + ext, nil
	// fileURL, err := url.Parse(*fullURL)

	// if err != nil {
	// 	log.Fatalf("Error occurred when parsing URL: %q. Error: %v", *fullURL, err)
	// }

	// path := fileURL.Path
	// filename := strings.ReplaceAll(fileURL.Host, "/", "_") + "_" + strings.ReplaceAll(path, "/", "_")

	// return &filename
	// return fullURL
}

func createFile(filename *string) *os.File {
	file, err := os.Create(*filename)

	if err != nil {
		log.Fatalf("Error occurred when creating image file: %q. Error: %v", *filename, err)
	}

	return file
}

func createDirIfNotExists(dirPath string) error {
	_, fileStatErr := os.Stat(dirPath)
	if fileStatErr == nil {
		log.Printf("Directory %q already exists. Skipping.\n", dirPath)
	}

	err := os.MkdirAll(dirPath, 0777)

	if err != nil {
		log.Printf("Error occured creating directory: %q. Error: %v\n", dirPath, err)
	}

	return err
}

func writeJSON(filename *string, v interface{}) error {
	file := createFile(filename)

	encoder := json.NewEncoder(file)

	err := encoder.Encode(v)

	return err
}
