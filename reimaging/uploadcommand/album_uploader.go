package uploadcommand

import (
	"fmt"
	"github.com/forgoty/go-reimaging/reimaging/validator"
	vkw "github.com/forgoty/go-reimaging/reimaging/vkwrapper"
	"os"
)

type result struct {
	err error
}

var extensions = map[string]bool{
	"jpg": true,
	"png": true,
	"gif": true,
	"bmp": true,
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

func (au *AlbumUploader) Upload(albumId int, filepath string) {
	files := validator.ReadDir(filepath)
	//semaphoreChan := make(chan struct{}, 1)
	// resultsChan := make(chan *result)
	defer func() {
		// close(semaphoreChan)
		// close(resultsChan)
	}()
	for _, file := range files {
		res := au.uploadPhoto(albumId, file)
		fmt.Println(res.err)
	}
	// var results []result
	// fileLen := len(files)
	// for {
	// 	result := <-resultsChan
	// 	fmt.Println(result.err)
	// 	results = append(results, *result)
	// 	if len(results) == fileLen {
	// 		break
	// 	}
	// }
}

func (au *AlbumUploader) uploadPhoto(albumId int, filepath string) *result {
	// defer func() {
	// 	<-semaphoreChan
	// }()
	res := &result{}
	file, err := os.Open(filepath)
	defer file.Close()
	if err != nil {
		res = &result{err}
	} else {
		err = au.vkWrapper.UploadPhoto(albumId, file)
		res = &result{err}
	}
	return res
	// resultsChan <- res
}
