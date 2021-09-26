package uploadcommand

import (
	"fmt"
	"github.com/forgoty/go-reimaging/pkg/reimaging/progressbar"
	"github.com/forgoty/go-reimaging/pkg/reimaging/validator"
	vkw "github.com/forgoty/go-reimaging/pkg/reimaging/vkwrapper"
	"net/http"
	"os"
)

var FilesInPostRequest int = 5

type AlbumUploader struct {
	vkWrapper vkw.VKWrapper
}

func NewAlbumUploader(vkWrapper vkw.VKWrapper) *AlbumUploader {
	return &AlbumUploader{
		vkWrapper: vkWrapper,
	}
}

func (au *AlbumUploader) Upload(albumId int, filepath, title string) {
	files := validator.ReadDir(filepath)
	lenFiles := len(files)
	uploadServer := au.vkWrapper.GetUploadServer(albumId)
	fileGroups := createFileGroups(files)
	client := http.Client{}
	semaphoreChan := getSemaphoreChannel(lenFiles)
	errCh := make(chan error)
	defer func() {
		close(semaphoreChan)
		close(errCh)
	}()

	bar := progressbar.NewProgressBar(lenFiles, title)
	for _, group := range fileGroups {
		groupUploader := NewFileGroupUploader(&client, au.vkWrapper, semaphoreChan, errCh)
		go groupUploader.Upload(group, uploadServer, albumId)
	}
	waitUpload(errCh, bar, len(fileGroups))

}

func getSemaphoreChannel(filesLen int) chan struct{} {
	if filesLen < 490 {
		return make(chan struct{}, 25)
	}
	return make(chan struct{}, 3)
}

func createFileGroups(filePaths []string) [][]string {
	var ret [][]string
	min := func(a, b int) int {
		if a <= b {
			return a
		}
		return b
	}

	for i := 0; i < len(filePaths); i += FilesInPostRequest {
		ret = append(ret, filePaths[i:min(i+FilesInPostRequest, len(filePaths))])
	}
	return ret

}

func waitUpload(errCh chan error, bar progressbar.ProgressBarHandler, total int) {
	var results []error
	var errors []error
	for {
		err := <-errCh
		if err != nil {
			checkErrors(bar, errors, err)
		}
		err = bar.Add(FilesInPostRequest)
		if err != nil {
			checkErrors(bar, errors, err)
		}
		results = append(results, err)
		if len(results) == total {
			err = bar.Finish()
			if err != nil {
				fmt.Println(err)
			}
			break
		}
	}
}

func checkErrors(bar progressbar.ProgressBarHandler, errors []error, err error) {
	errors = append(errors, err)
	if len(errors) > 3 {
		err := bar.Finish()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println()
		fmt.Println("Too many errors occured recently")
		os.Exit(1)
	}
}
