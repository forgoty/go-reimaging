package cmd

import (
	"fmt"
	vkapi "github.com/SevereCloud/vksdk/5.92/api"
	"os"
)

func GetAlbums(vk *vkapi.VK, userId string) vkapi.PhotosGetAlbumsResponse {
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
	return response
}
