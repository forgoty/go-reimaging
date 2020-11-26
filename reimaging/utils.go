package reimaging

import (
	"fmt"
	"os"
	"strconv"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
)

func GetAlbums(vk *api.VK, userId int) []object.PhotosPhotoAlbumFull {
	response_params := params.NewPhotosGetAlbumsBuilder()
	response_params.OwnerID(userId)
	response_params.NeedSystem(System)
	response, vkErr := vk.PhotosGetAlbums(response_params.Params)
	if vkErr != nil {
		fmt.Println(vkErr)
		os.Exit(1)
	}
	return response.Items
}

func GetPhotoUrls(vk *api.VK, userId, albumId, offsetInt int) []string {
	response_params := params.NewPhotosGetBuilder()
	response_params.OwnerID(userId)
	response_params.AlbumID(strconv.Itoa(albumId))
	response_params.Offset(offsetInt)
	response_params.PhotoSizes(true)
	response_params.Count(1000)
	response, vkErr := vk.PhotosGet(response_params.Params)
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

func GetVk() *api.VK {
	token := os.Getenv("VK_TOKEN")
	vk := api.NewVK(token)
	return vk
}
