package downloadcommand

import (
	"fmt"
	"os"
	"strconv"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
)

type VKWrapper struct {
	vk *api.VK
}

func NewVKWrapper() *VKWrapper {
	return &VKWrapper{getVk()}
}

func (vkw *VKWrapper) GetPhotoURLs(album PhotoAlbum, offset int) []string {
	response_params := params.NewPhotosGetBuilder()
	response_params.OwnerID(album.OwnerID)
	response_params.AlbumID(strconv.Itoa(album.ID))
	response_params.Offset(offset)
	response_params.PhotoSizes(true)
	response_params.Count(1000)
	response, vkErr := vkw.vk.PhotosGet(response_params.Params)
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

func (vkw *VKWrapper) GetAlbums(userID int, NeedSystem bool) []PhotoAlbum {
	response_params := params.NewPhotosGetAlbumsBuilder()
	response_params.OwnerID(userID)
	response_params.NeedSystem(NeedSystem)
	response, vkErr := vkw.vk.PhotosGetAlbums(response_params.Params)
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

func (vkw *VKWrapper) CreateAlbum(title string) PhotoAlbum{
	response_params := params.NewPhotosCreateAlbumBuilder()
	response_params.Title(title)
	response_params.PrivacyView([]string{"only_me"})
	rawAlbum, vkErr := vkw.vk.PhotosCreateAlbum(response_params.Params)
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

func getVk() *api.VK {
	return api.NewVK(os.Getenv("VK_TOKEN"))
}

type PhotoAlbum struct {
	ID          int
	OwnerID     int
	Size        int
	Title       string
}
