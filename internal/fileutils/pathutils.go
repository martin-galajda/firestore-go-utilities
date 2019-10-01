package fileutils

import (
	"path/filepath"
)

// getFilepathExtension extracts extension
// from the full filepath.
func getFilepathExtension(pathToFile string) string {
	ext := filepath.Ext(pathToFile)

	return ext
}

// GetFilenameFromFilepath extracts filename
// from the full filepath.
func GetFilenameFromFilepath(pathToFile string) string {
	_, filename := filepath.Split(pathToFile)

	return filename
}

