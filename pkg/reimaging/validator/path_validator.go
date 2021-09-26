package validator

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func ValidateDownloadDir(path string) (string, error) {
	path, error := validatePath(path)
	return path, error
}

func validatePath(path string) (string, error) {
	path = getAbsPath(path)

	if res, err := isDir(path); err != nil || !res {
		return "", errors.New("Error: Giving Path is not valid")
	}

	return path, nil

}

func getAbsPath(path string) string {
	if path == "" {
		path, _ = os.Getwd()
		return path
	}
	path, _ = filepath.Abs(path)
	return path
}

func isDir(path string) (bool, error) {
	file, err := os.Stat(path)
	if os.IsNotExist(err) || !file.IsDir() {
		return false, err
	}
	return true, nil
}

func ReadDir(root string) []string {
	var extensions = [4]string{"jpg", "bmp", "png", "gif"}
	var files []string
	fileInfo, _ := ioutil.ReadDir(root)
	for _, file := range fileInfo {
		for _, ext := range extensions {
			if strings.HasSuffix(strings.ToLower(file.Name()), ext) {
				files = append(files, root+string(os.PathSeparator)+file.Name())
			}
		}
	}
	return files
}
