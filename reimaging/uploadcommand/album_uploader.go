package uploadcommand

import (
	vkw "github.com/forgoty/go-reimaging/reimaging/vkwrapper"
)

type AlbumUploader struct {
	vkWrapper vkw.VKWrapper
	UploadPath string
}

func NewAlbumUploader(vk vkw.VKWrapper, uploadPath string) *AlbumUploader {
	return &AlbumUploader{
		vkWrapper: vk,
		UploadPath: uploadPath,
	}
}

func (au *AlbumUploader) CreateAlbum(title string) vkw.PhotoAlbum {
	return au.vkWrapper.CreateAlbum(title)
}

func (au *AlbumUploader) getUploadServer(album vkw.PhotoAlbum) string {
	return au.vkWrapper.GetUploadServer(album.ID)
}
