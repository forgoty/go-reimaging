package vkwrapper

import (
	"github.com/stretchr/testify/mock"
)

type MockVKWrapper struct{
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
	return PhotoAlbum{
		ID: 12345,
		OwnerID: 12345,
		Size: 1,
		Title: "mock",
	}
}
