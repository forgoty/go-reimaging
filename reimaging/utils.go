package reimaging

import (
	"fmt"
	"os"
	"strconv"
)

type PhotoSize struct {
	Height int
	URL    string
	Type   string
	Width  int
}

func GetAlbums(vk *vkapi.VK, userId string) []object.PhotosPhotoAlbumFull {
	var needSystem string
	if System {
		needSystem = "1"
	} else {
		needSystem = "0"
	}
	params := map[string]string{
		"need_system": needSystem,
		"owner_id":    userId,
	}
	response, vkErr := vk.PhotosGetAlbums(params)
	if vkErr.Code != 0 {
		fmt.Println(vkErr.Message)
		os.Exit(1)
	}
	return response.Items
}

func GetPhotoUrls(vk *vkapi.VK, userId, albumId, offsetInt int) []string {
	userID := strconv.Itoa(userId)
	albumID := strconv.Itoa(albumId)
	offset := strconv.Itoa(offsetInt)
	params := map[string]string{
		"album_id":     albumID,
		"owner_id":     userID,
		"offset":       offset,
		"photos_sizes": "1",
		"count":        "1000",
	}
	response, vkErr := vk.PhotosGet(params)
	if vkErr.Code != 0 {
		fmt.Println(vkErr.Message)
	}

	urls := []string{}
	for _, photo := range response.Items {
		urls = append(urls, PhotoSize(photo.Sizes[len(photo.Sizes)-1]).URL)
	}
	return urls
}
