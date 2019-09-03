package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"encoding/json"
)

func downloadFile(client *http.Client, url *string, pathToOutputDir string) {
	resp, err := client.Get(*url)
	if err != nil {
		log.Fatalf("Error occurred when downloading image file from URL: %q. Error: %v", *url, err)
	}
	defer resp.Body.Close()

	filename := buildFilename(url)
	fullPath := path.Join(pathToOutputDir, *filename)

	file := createFile(&fullPath)

	size, err := io.Copy(file, resp.Body)
	defer file.Close()

	log.Printf("Downloaded a file %q with size %d\n", *filename, size)
}

func buildFilename(fullUrl *string) *string {
	fileUrl, err := url.Parse(*fullUrl)

	if err != nil {
		log.Fatalf("Error occurred when parsing URL: %q. Error: %v", *fullUrl, err)
	}

	path := fileUrl.Path
	filename := strings.ReplaceAll(fileUrl.Host, "/", "_") + "_" + strings.ReplaceAll(path, "/", "_")

	return &filename
}

func createFile(filename *string) *os.File {
	file, err := os.Create(*filename)

	if err != nil {
		log.Fatalf("Error occurred when creating image file: %q. Error: %v", *filename, err)
	}

	return file
}

func writeJSON(filename *string, v interface{}) error {
	file := createFile(filename)

	encoder := json.NewEncoder(file)

	err := encoder.Encode(v)

	return err
}
