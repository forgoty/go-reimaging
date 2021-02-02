package downloadcommand

import (
	"testing"

	vkw "github.com/forgoty/go-reimaging/reimaging/vkwrapper"
)

var testAlbum = vkw.PhotoAlbum{
	ID: 12345,
	OwnerID: 12345,
	Size: 1,
	Title: "mock",
}

func TestAlbumDownloaderDownloadAll(t *testing.T) {
	mockedVKWrapper := new(vkw.MockVKWrapper)
	mockedVKWrapper.On("GetAlbums", 12345, false).Return([]vkw.PhotoAlbum{})

	options := NewDownloadOptions(12345, "path", false)

	albumDownloader := NewAlbumDownloader(mockedVKWrapper, options)
	albumDownloader.DownloadAll()

	mockedVKWrapper.AssertExpectations(t)
}

func TestAlbumDownloaderDownloadByID(t *testing.T) {
	mockedVKWrapper := new(vkw.MockVKWrapper)
	mockedVKWrapper.On("GetAlbums", 12345, false).Return([]vkw.PhotoAlbum{})

	options := NewDownloadOptions(12345, "path", false)

	albumDownloader := NewAlbumDownloader(mockedVKWrapper, options)
	albumDownloader.DownloadAlbumByID(12345)

	mockedVKWrapper.AssertExpectations(t)
}

func mockedDownloadPhoto(
	metaImage metaDownloadedPhoto,
	done chan bool,
	errch chan bool,
) {}

