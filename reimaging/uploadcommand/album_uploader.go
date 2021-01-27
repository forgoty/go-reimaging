package uploadcommand

import (
	vkw "github.com/forgoty/go-reimaging/reimaging/vkwrapper"
)

type AlbumUploader struct {
	VKWrapper *vkw.VKWrapper
	UploadPath string
}

func NewAlbumUploader(uploadPath string) *AlbumUploader {
	return &AlbumUploader{
		VKWrapper: vkw.NewVKWrapper(),
		UploadPath: uploadPath,
	}
}
