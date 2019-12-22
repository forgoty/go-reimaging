package cmd

import (
	"fmt"
	"os"
	vkapi "github.com/SevereCloud/vksdk/5.92/api"
	object "github.com/SevereCloud/vksdk/5.92/object"
)

func GetAlbums(vk *vkapi.VK, userId string) []object.PhotosPhotoAlbumFull  {
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

func GetPhotos(vk *vkapi.VK, userId, albumId, offset string) []object.PhotosPhoto {
	params := map[string]string{
		"album_id": albumId,
		"owner_id":    userId,
		"offset":    offset,
		"photos_sizes":    "1",
		"count":    "1000",
	}
	response, vkErr := vk.PhotosGet(params)
	if vkErr.Code != 0 {
		fmt.Println(vkErr.Message)
	}
	return response.Items
}
