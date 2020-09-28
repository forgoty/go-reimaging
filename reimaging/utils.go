package reimaging

import (
	"fmt"
	"os"
	"strconv"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
)

type PhotoSize struct {
	Height int
	URL    string
	Type   string
	Width  int
}

func GetAlbums(vk *api.VK, userId string) []object.PhotosPhotoAlbumFull {
	return []object.PhotosPhotoAlbumFull{}
	// var needSystem string
	// if System {
	// 	needSystem = "1"
	// } else {
	// 	needSystem = "0"
	// }
	// params := map[string]string{
	// 	"need_system": needSystem,
	// 	"owner_id":    userId,
	// }
	// response, vkErr := vk.PhotosGetAlbums(params)
	// if vkErr.Code != 0 {
	// 	fmt.Println(vkErr.Message)
	// 	os.Exit(1)
	// }
	// return response.Items
}

func GetPhotoUrls(vk *api.VK, userId, albumId, offsetInt int) []string {
	photo_params := params.NewPhotosGetBuilder()
	photo_params.Count(1000)
	photo_params.AlbumID(strconv.Itoa(albumId))
	photo_params.Offset(offsetInt)
	photo_params.OwnerID(userId)
	photo_params.PhotoSizes(true)
	// userID := strconv.Itoa(userId)
	// albumID := strconv.Itoa(albumId)
	// offset := strconv.Itoa(offsetInt)
	// params := map[string]string{
	// 	"album_id":     albumID,
	// 	"owner_id":     userID,
	// 	"offset":       offset,
	// 	"photos_sizes": "1",
	// 	"count":        "1000",
	// }
	response, _ := vk.PhotosGet(photo_params.Params)
	// if vkErr.Code != 0 {
	// 	fmt.Println(vkErr.Message)
	// }

	urls := []string{}
	for _, photo := range response.Items {
		urls = append(urls, photo.Sizes[len(photo.Sizes)-1])
	}
	return urls
}

func GetVk() *api.VK {
	token := os.Getenv("VK_TOKEN")
	vk := api.NewVK(token)
	return vk
}
