package uploadcommand

import (
	vkw "github.com/forgoty/go-reimaging/reimaging/vkwrapper"
)

type AlbumUploader struct {
	vkWrapper vkw.VKWrapper
	UploadPath string
}

func NewAlbumUploader(uploadPath string) *AlbumUploader {
	return &AlbumUploader{
		vkWrapper: vkw.NewVKWrapper(),
		UploadPath: uploadPath,
	}
}
