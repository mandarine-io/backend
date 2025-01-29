package file

import (
	"os"
	"path/filepath"
)

func GetFileNameWithoutExt(filePath string) string {
	return filePath[:len(filePath)-len(filepath.Ext(filePath))]
}

func GetFilesFromDir(dirPath string) ([]string, error) {
	dir, err := os.Open(filepath.Clean(dirPath))
	if err != nil {
		return make([]string, 0), err
	}

	files, err := dir.Readdir(0)
	if err != nil {
		return make([]string, 0), err
	}

	fileNames := make([]string, 0)
	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, file.Name())
		}
	}

	return fileNames, nil
}
