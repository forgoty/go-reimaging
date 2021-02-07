package uploadcommand

import (
	"testing"

	vkw "github.com/forgoty/go-reimaging/reimaging/vkwrapper"
)

func TestAlbumCreateAlbum(t *testing.T) {
	mockedVKWrapper := new(vkw.MockVKWrapper)
	mockedVKWrapper.On("CreateAlbum", vkw.TestName).Return(vkw.TestAlbum)


	albumUploader := NewAlbumUploader(mockedVKWrapper, vkw.TestName)
	albumUploader.CreateAlbum(vkw.TestName)

	mockedVKWrapper.AssertExpectations(t)
}

func TestGetUploadServer(t *testing.T) {
	mockedVKWrapper := new(vkw.MockVKWrapper)
	mockedVKWrapper.On("GetUploadServer", vkw.TestID).Return(vkw.TestName)


	albumUploader := NewAlbumUploader(mockedVKWrapper, vkw.TestName)
	albumUploader.getUploadServer(vkw.TestAlbum)

	mockedVKWrapper.AssertExpectations(t)

}
