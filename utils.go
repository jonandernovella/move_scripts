package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func formatBytes(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

func getAbsoluteDirectory(path string) (string, error) {
	isAbsolute := filepath.IsAbs(path)
	err := error(nil)

	if !isAbsolute {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return "", errors.New("Error converting path to absolute: " + err.Error())
		}
		path = absPath
	}

	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return "", errors.New("Path does not exist: " + path)
	}

	if !fileInfo.IsDir() {
		return "", errors.New("Path is not a directory: " + path)
	}
	return path, nil
}
