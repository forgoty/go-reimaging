package vkwrapper

import (
	"fmt"
	"os"
	"strconv"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
)

type VKWrapper interface {
	GetPhotoURLs(album PhotoAlbum, offset int) []string
	GetAlbums(userID int, NeedSystem bool) []PhotoAlbum
	CreateAlbum(title string) PhotoAlbum
	GetUploadServer(id int) string
}

type VkAPIWrapper struct {
	vk *api.VK
}

func NewVKWrapper() *VkAPIWrapper {
	return &VkAPIWrapper{getVk()}
}

func (vkw *VkAPIWrapper) GetPhotoURLs(album PhotoAlbum, offset int) []string {
	responseParams := params.NewPhotosGetBuilder()
	responseParams.OwnerID(album.OwnerID)
	responseParams.AlbumID(strconv.Itoa(album.ID))
	responseParams.Offset(offset)
	responseParams.PhotoSizes(true)
	responseParams.Count(1000)
	response, vkErr := vkw.vk.PhotosGet(responseParams.Params)
	if vkErr != nil {
		fmt.Println(vkErr)
	}

	urls := []string{}
	for _, photo := range response.Items {
		if len(photo.Sizes) > 0 {
			// the most high-rez picture is always last
			urls = append(urls, photo.Sizes[len(photo.Sizes)-1].BaseImage.URL)
		}
	}
	return urls
}

func (vkw *VkAPIWrapper) GetAlbums(userID int, NeedSystem bool) []PhotoAlbum {
	responseParams := params.NewPhotosGetAlbumsBuilder()
	responseParams.OwnerID(userID)
	responseParams.NeedSystem(NeedSystem)
	response, vkErr := vkw.vk.PhotosGetAlbums(responseParams.Params)
	if vkErr != nil {
		fmt.Println(vkErr)
		os.Exit(1)
	}
	albums := []PhotoAlbum{}
	for _, rawAlbum := range response.Items {
		albums = append(
			albums,
			PhotoAlbum{
				ID: rawAlbum.ID,
				OwnerID: rawAlbum.OwnerID,
				Size: rawAlbum.Size,
				Title: rawAlbum.Title,
			},
		)
	}

	return albums
}

func (vkw *VkAPIWrapper) CreateAlbum(title string) PhotoAlbum{
	responseParams := params.NewPhotosCreateAlbumBuilder()
	responseParams.Title(title)
	responseParams.PrivacyView([]string{"only_me"})
	rawAlbum, vkErr := vkw.vk.PhotosCreateAlbum(responseParams.Params)
	if vkErr != nil {
		fmt.Println(vkErr)
		os.Exit(1)
	}

	return PhotoAlbum{
		ID: rawAlbum.ID,
		OwnerID: rawAlbum.OwnerID,
		Size: rawAlbum.Size,
		Title: rawAlbum.Title,
	}
}

func (vkw *VkAPIWrapper) GetUploadServer(id int) string {
	responseParams := params.NewPhotosGetUploadServerBuilder()
	responseParams.AlbumID(id)
	uploadServerResponse, vkErr := vkw.vk.PhotosGetUploadServer(responseParams.Params)
	if vkErr != nil {
		fmt.Println(vkErr)
		os.Exit(1)
	}
	return uploadServerResponse.UploadURL
}

func getVk() *api.VK {
	return api.NewVK(os.Getenv("VK_TOKEN"))
}

type PhotoAlbum struct {
	ID          int
	OwnerID     int
	Size        int
	Title       string
}
