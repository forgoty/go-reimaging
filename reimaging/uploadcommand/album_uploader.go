package uploadcommand

import (
	"fmt"
	"github.com/forgoty/go-reimaging/reimaging/validator"
	vkw "github.com/forgoty/go-reimaging/reimaging/vkwrapper"
	progressbar "github.com/schollz/progressbar/v3"
	"net/http"
	"os"
)

var FilesInPostRequest int = 5

type result struct {
	err error
}

type AlbumUploader struct {
	vkWrapper vkw.VKWrapper
}

func NewAlbumUploader() *AlbumUploader {
	return &AlbumUploader{
		vkWrapper: vkw.NewVKWrapper(),
	}
}

func (au *AlbumUploader) CreateAlbum(title string) vkw.PhotoAlbum {
	return au.vkWrapper.CreateAlbum(title)
}

func (au *AlbumUploader) GetAlbumsByIDs(ids []int) []vkw.PhotoAlbum {
	return au.vkWrapper.GetAlbumsByAlbumIds(ids)
}

func (au *AlbumUploader) Upload(albumId int, filepath, title string) {
	files := validator.ReadDir(filepath)
	lenFiles := len(files)
	uploadServer := au.vkWrapper.GetUploadServer(albumId)
	fileGroups := au.createFileGroups(files)
	client := http.Client{}
	semaphoreChan := make(chan struct{}, calculateSemaphoreCount(lenFiles))
	errCh := make(chan error)
	var results []error

	bar := getProgressBar(int64(lenFiles), title)
	for _, group := range fileGroups {
		go au.vkWrapper.UploadFileGroup(&client, group, uploadServer, albumId, semaphoreChan, errCh)
	}
	for {
		err := <-errCh
		if err != nil {
			fmt.Println(err)
		} else {
			state := bar.State().CurrentBytes
			bar.Add(calculateAdd(int(state), lenFiles))
		}
		results = append(results, err)
		if len(filterNotNils(results)) > 3 {
			bar.Finish()
			fmt.Println()
			fmt.Println("Too many errors occured recently")
			os.Exit(1)
		}
		if len(results) == len(fileGroups) {
			bar.Finish()
			break
		}
	}
}

func calculateAdd(current, max int) int {
	if current+FilesInPostRequest >= max {
		return max - current
	}
	return FilesInPostRequest
}

func calculateSemaphoreCount(filesLen int) int {
	if filesLen < 490 {
		return 25
	}
	return 3
}

func filterNotNils(results []error) []error {
	var res []error
	for _, i := range results {
		if i != nil {
			res = append(res, i)
		}
	}
	return res
}

func (au *AlbumUploader) createFileGroups(filePaths []string) [][]string {
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

func getProgressBar(max int64, title string) *progressbar.ProgressBar {
	theme := progressbar.Theme{
		Saucer:        "=",
		SaucerHead:    ">",
		SaucerPadding: " ",
		BarStart:      "[",
		BarEnd:        "]",
	}
	return progressbar.NewOptions64(
		max,
		progressbar.OptionShowIts(),
		progressbar.OptionSetItsString("photos"),
		progressbar.OptionSetDescription(title),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionShowCount(),
		progressbar.OptionSetTheme(theme),
	)

}
