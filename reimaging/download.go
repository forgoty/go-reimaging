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
	"sync"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/object"
	progressbar "github.com/schollz/progressbar/v3"
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
	userId, _ := strconv.Atoi(args[0])

	_, error := validators.ValidateDownloadDir(path)
	if error != nil {
		fmt.Println(error)
		os.Exit(1)
	}

	vk := GetVk()

	if AlbumId != 0 {
		downloadAlbumById(vk, userId, AlbumId)
	} else {
		downloadAlbums(vk, userId)
	}
}

func downloadAlbumById(vk *api.VK, userId, albumId int) {
	albums := GetAlbums(vk, userId)
	for _, album := range albums {
		if album.ID == albumId {
			downloadAlbum(vk, album)
		}
	}
}

func downloadAlbums(vk *api.VK, userId int) {
	var wg sync.WaitGroup

	albums := GetAlbums(vk, userId)
	for _, album := range albums {
		wg.Add(1)
		go downloadAlbumWaitGroup(&wg, vk, album)
	}
	wg.Wait()
}

func downloadAlbumWaitGroup(wg *sync.WaitGroup, vk *api.VK, album object.PhotosPhotoAlbumFull) {
	defer wg.Done()
	downloadAlbum(vk, album)
}

func downloadAlbum(vk *api.VK, album object.PhotosPhotoAlbumFull) {
	offsets := getOffset(album.Size)
	photosUrls := []string{}
	for _, offset := range offsets {
		photosUrls = append(photosUrls, GetPhotoUrls(vk, album.OwnerID, album.ID, offset)...)
	}
	if len(photosUrls) > 0 {
		pathDir := createAlbumDir(album.Title)
		downloadPhotos(photosUrls, pathDir, album.Title)
	}
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

type metaDownloadedPhoto struct {
	url      string
	filename string
	path     string
}

func downloadPhotos(photosUrls []string, pathDir, albumTitle string) {
	photosCount := len(photosUrls)
	done := make(chan bool, photosCount)
	errch := make(chan error, photosCount)
	for _, url := range photosUrls {
		metaImage := getPhotoMetaData(url, pathDir)
		go downloadPhoto(metaImage, done, errch)
	}
	bar := progressbar.Default(int64(photosCount), albumTitle)
	for i := 0; i < photosCount; i++ {
		if ok := <-done; !ok {
			if err := <-errch; err != nil {
				fmt.Println(err)
			}
		}
		bar.Add(1)
	}
	close(done)
	close(errch)
}

func getPhotoMetaData(url, pathDir string) metaDownloadedPhoto {
	split := strings.Split(url, "/")
	filename := strings.ReplaceAll(split[len(split)-1], "-", "")
	path := pathDir + "/" + filename
	return metaDownloadedPhoto{url: url, filename: filename, path: path}
}

func downloadPhoto(metaImage metaDownloadedPhoto, done chan bool, errch chan error) {
	_, err := os.Stat(metaImage.path)
	if err == nil {
		// file exists
		done <- true
		errch <- nil
		return
	}

	response, err := http.Get(metaImage.url)
	if err != nil {
		errch <- err
		done <- false
		return
	}
	defer response.Body.Close()

	output, err := os.Create(metaImage.path)
	if err != nil {
		errch <- err
		done <- false
		return
	}
	defer output.Close()

	_, err = io.Copy(output, response.Body)
	if err != nil {
		errch <- err
		done <- false
		return
	}
	done <- true
	errch <- nil
	return
}
