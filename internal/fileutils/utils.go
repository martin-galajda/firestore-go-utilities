package fileutils

import (
	"fmt"
	"strings"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"io/ioutil"

	"encoding/json"
)

func DownloadAndSaveFile(client *http.Client, url *string, fileId string, pathToOutputDir string) {
	resp, err := client.Get(*url)
	if err != nil {
		log.Fatalf("Error occurred when downloading image file from URL: %q. Error: %v", *url, err)
	}
	defer resp.Body.Close()

	filename, err := BuildFilename(url, fileId)

	if err != nil {
		log.Printf("Failed to export file with url: %q\n. Skipping...", *url)
		return
	}

	fullPath := path.Join(pathToOutputDir, filename)

	file := CreateFile(&fullPath)

	size, err := io.Copy(file, resp.Body)
	defer file.Close()

	log.Printf("Downloaded and saved a file %q with size %d\n", filename, size)
}

func BuildFilename(fullURL *string, fileID string) (string, error) {
	// get rid of query string from URL
	*fullURL = regexp.MustCompile(`\?.+$`).ReplaceAllString(*fullURL, "")
	ext := getFilepathExtension(*fullURL)

	if ext == "" {
		errMsg := fmt.Sprintf("Couldn't extract extension from URL: %q. Skipping import...", *fullURL)
		err := fmt.Errorf(errMsg)
		log.Printf("Couldn't extract extension from URL: %q. Skipping import...", *fullURL)

		return "", err
	}


	return fileID + ext, nil
}

// CreateFile creates file in case it does not exist.
// It terminates program in case anything goes wrong.
func CreateFile(filename *string) *os.File {
	file, err := os.Create(*filename)

	if err != nil {
		log.Fatalf("Error occurred when creating image file: %q. Error: %v", *filename, err)
	}

	return file
}

// CreateDirIfNotExists creates directory
// in case it does not exist.
func CreateDirIfNotExists(dirPath string) error {
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

// WriteJSON encodes any value as JSON
// and saves it to the file specified by filename.
func WriteJSON(filename *string, v interface{}) error {
	file := CreateFile(filename)

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")

	err := encoder.Encode(v)

	return err
}

// ReadFileAsText reads content of the file into string.
func ReadFileAsText(pathToFile string) (string, error) {
	bytes, err := ioutil.ReadFile(pathToFile)

	if err != nil {
		return "", err
	}

	fileContent := string(bytes)

	return fileContent, nil
}

// GetFilePathsInDirWithExt gets slice of absolute path to files
// contained in the given directory 
// and that have matching extension that is provided to the function.
func GetFilePathsInDirWithExt(pathToDir, ext string) ([]string, error) {
	filesInDir, err := ioutil.ReadDir(pathToDir)

	if err != nil {
		return nil, err
	}

	results := make([]string, 0, len(filesInDir))

	for _, fileInfo := range filesInDir {
		if strings.Contains(fileInfo.Name(), ext) {
			results = append(results, path.Join(pathToDir, fileInfo.Name()))
		}
	}

	return results, nil
}
