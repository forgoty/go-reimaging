package cmd

import (
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"

	vkapi "github.com/SevereCloud/vksdk/5.92/api"
	object "github.com/SevereCloud/vksdk/5.92/object"
	"github.com/spf13/cobra"

	"github.com/forgoty/go-reimaging/cmd/auth"
	"github.com/forgoty/go-reimaging/cmd/validators"
)

var downloadCmd = &cobra.Command{
	Use:   "download USERID",
	Short: "Download photo albums of specific user",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("More then one argument provided")
		}

		_, error := strconv.Atoi(args[0])
		if error != nil {
			return errors.New("Unvalid USERID has been provided. Need Integer")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		download(args)
	},
}

var path string

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().BoolVarP(&Auth, "auth", "a", false, "Enable authorization")
	downloadCmd.Flags().BoolVarP(&System, "system", "s", false, "Enable system albums for download")
	downloadCmd.Flags().IntVarP(&AlbumId, "album-id", "", 0, "Use specific album ID to download")
	downloadCmd.Flags().StringVarP(&path, "path", "p", "", "Set Download Folder")
}

func download(args []string) {
	userId := args[0]

	_, error := validators.ValidateDownloadDir(path)
	if error != nil {
		fmt.Println(error)
		os.Exit(1)
	}
	vk := auth.GetClient(Auth)
	albums := GetAlbums(vk, userId)
	for _, album := range albums {
		downloadAlbum(vk, album)
	}
}

func downloadAlbum(vk *vkapi.VK, album object.PhotosPhotoAlbumFull) error {
	createAlbumDir(album.Title)
	offsets := getOffset(album.Size)
	photos := []string{}
	for _, offset := range offsets {
		photos = append(photos, GetPhotoUrls(vk, album.OwnerID, album.ID, offset)...)
	}
	return nil
}

func createAlbumDir(title string) {
	pathDir := filepath.Join(path, title)
	if _, serr := os.Stat(pathDir); serr != nil {
		os.MkdirAll(pathDir, os.ModePerm)
	}
}

func getOffset(size int) []int {
	var maxCount int = 1000
	d := float64(size) / float64(maxCount)
	offset := int(math.Ceil(d))

	offsets := make([]int, offset)

	buf := 0
	for i := 0; i < offset; i++ {
		offsets[i] = buf
		buf += maxCount
	}
	return offsets
}
