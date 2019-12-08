package validators

import (
	"errors"
	"os"
	"path/filepath"
)

func ValidateDownloadDir(path string) (string, error) {
	path, error := validatePath(path)
	return path, error
}

func validatePath(path string) (string, error) {
	path = getAbsPath(path)

	if res, err := isDir(path); err != nil || res != true {
		return "", errors.New("Giving Path is not valid")
	}

	return path, nil

}

func getAbsPath(path string) string {
	if path == "" {
		path, _ = os.Getwd()
	} else {
		path, _ = filepath.Abs(path)
	}
	return path
}

func isDir(path string) (bool, error) {
	file, err := os.Stat(path)
	if os.IsNotExist(err) || !file.IsDir() {
		return false, err
	}
	return true, nil
}

// For Upload Command

// func IsPathContainsImages(path string) bool {
// 	files := IOReadDir(path)
// 	if len(files) > 0 {
// 		return true
// 	} else {
// 		return false
// 	}
// }

// func IOReadDir(root string) []string {
// 	var extensions = [4]string{"jpg", "bmp", "png", "gif"}
// 	var files []string
// 	fileInfo, _ := ioutil.ReadDir(root)
// 	for _, file := range fileInfo {
// 		for _, ext := range extensions {
// 			if strings.HasSuffix(file.Name(), ext) {
// 				files = append(files, file.Name())
// 			}
// 		}
// 	}
// 	return files
// }
