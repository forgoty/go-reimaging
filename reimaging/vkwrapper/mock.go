package vkwrapper

import (
	"github.com/stretchr/testify/mock"
)

var (
	TestName  = "mock"
	TestID    = 123
	TestAlbum = PhotoAlbum{
		ID:      TestID,
		OwnerID: TestID,
		Size:    666,
		Title:   TestName,
	}
)

type MockVKWrapper struct {
	mock.Mock
}

func (m *MockVKWrapper) GetPhotoURLs(album PhotoAlbum, offset int) []string {
	args := m.Called(album, offset)
	return []string{args.String(0)}
}
func (m *MockVKWrapper) GetAlbums(userID int, NeedSystem bool) []PhotoAlbum {
	m.Called(userID, NeedSystem)
	return []PhotoAlbum{}
}

func (m *MockVKWrapper) CreateAlbum(title string) PhotoAlbum {
	m.Called(title)
	return TestAlbum
}

func (m *MockVKWrapper) GetUploadServer(id int) string {
	m.Called(id)
	return TestName
}

func (m *MockVKWrapper) GetAlbumsByAlbumIds(albumIDs []int) []PhotoAlbum {
	m.Called(albumIDs)
	return []PhotoAlbum{}
}

func (m *MockVKWrapper) PhotosSave(_ []byte, albumId int) error {
	m.Called(albumId)
	return nil
}
