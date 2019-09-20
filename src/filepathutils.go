package main

import (
	"path/filepath"
	"regexp"
)

func getFilepathExtension(pathToFile string) string {
	ext := filepath.Ext(pathToFile)

	return ext
}

func getFilenameFromFilepath(pathToFile string) string {
	_, filename := filepath.Split(pathToFile)

	return filename
}

func trimFilepathExtension(pathToFile string) string {
	ext := getFilepathExtension(pathToFile)

	re := regexp.MustCompile(ext + "$")

	return re.ReplaceAllString(pathToFile, "")
}

