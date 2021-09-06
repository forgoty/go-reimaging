package uploadcommand

import (
	"fmt"
	"github.com/forgoty/go-reimaging/reimaging/validator"
	vkw "github.com/forgoty/go-reimaging/reimaging/vkwrapper"
	progressbar "github.com/schollz/progressbar/v3"
	"net/http"
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
	uploadServer := au.vkWrapper.GetUploadServer(albumId)
	fileGroups := au.createFileGroups(files)
	client := http.Client{}
	bar := getProgressBar(int64(len(files)), title)
	for _, group := range fileGroups {
		err := au.vkWrapper.UploadFileGroup(&client, group, uploadServer, albumId)
		if err != nil {
			fmt.Println(err)
			continue
		}
		bar.Add(len(group))
	}
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
