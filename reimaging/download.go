package reimaging

import (
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/schollz/progressbar"
	"github.com/spf13/cobra"

	"github.com/forgoty/go-reimaging/reimaging/validators"
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

	albums := GetAlbums(vk, userId)
	for _, album := range albums {
		downloadAlbum(vk, album)
	}
}

func downloadAlbum(vk *vkapi.VK, album object.PhotosPhotoAlbumFull) error {
	pathDir := createAlbumDir(album.Title)
	offsets := getOffset(album.Size)
	photosUrls := []string{}
	for _, offset := range offsets {
		photosUrls = append(photosUrls, GetPhotoUrls(vk, album.OwnerID, album.ID, offset)...)
	}
	downloadPhotos(photosUrls, pathDir)
	return nil
}

func createAlbumDir(title string) string {
	pathDir := filepath.Join(path, title)
	if _, serr := os.Stat(pathDir); serr != nil {
		os.MkdirAll(pathDir, os.ModePerm)
	}
	return pathDir
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

func downloadPhotos(photosUrls []string, pathDir string) {
	done := make(chan bool)
	parallels := 0
	for _, url := range photosUrls {
		parallels++
		go downloadPhoto(url, pathDir, done)
	}
	bar := progressbar.New(parallels)
	for i := 0; i < parallels; i++ {
		<-done
		bar.Add(1)
	}
}

func downloadPhoto(url, pathDir string, done chan bool) {
	split := strings.Split(url, "/")
	fileName := strings.ReplaceAll(split[len(split)-1], "-", "")
	path := pathDir + "/" + fileName

	if _, err := os.Stat(path); err == nil {
		return
	}

	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error while downloading", url)
	}
	defer response.Body.Close()

	output, err := os.Create(path)
	if err != nil {
		fmt.Printf("Error while creating", fileName)
	}
	defer output.Close()

	_, err = io.Copy(output, response.Body)
	if err != nil {
		fmt.Printf("Error while saving", url)
	}
	done <- true
	return
}
