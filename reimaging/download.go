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
	"time"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/object"
	progressbar "github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
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

	_, error := validateDownloadDir(path)
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
	albums := GetAlbums(vk, userId)
	for _, album := range albums {
	    downloadAlbum(vk, album)
	    fmt.Println()
	}
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

	done := make(chan bool, 1)
	defer close(done)

	errch := make(chan error, 1)
	defer close(errch)

	for _, url := range photosUrls {
		metaImage := getPhotoMetaData(url, pathDir)
		go downloadPhoto(metaImage, done, errch)
	}

	bar := getProgressBar(int64(photosCount), albumTitle)
	for i := 0; i < photosCount; i++ {
		if ok := <-done; !ok {
			if err := <-errch; err != nil {
				fmt.Println(err)
			}
		} else {
			<-errch
		}
		bar.Add(1)
	}
}

func getPhotoMetaData(url, pathDir string) metaDownloadedPhoto {
	filename := getFileName(url)
	path := pathDir + "/" + filename
	return metaDownloadedPhoto{url: url, filename: filename, path: path}
}

func getFileName(url string) string{
	split := strings.Split(url, "/")
	filename := strings.ReplaceAll(split[len(split)-1], "-", "")
	if index := strings.Index(filename, "?"); index > 0 {
		filename = filename[0:index]
	}
	return filename
}

func downloadPhoto(metaImage metaDownloadedPhoto, done chan bool, errch chan error) {
	_, err := os.Stat(metaImage.path)
	if err == nil {
		// file exists
		done <- true
		errch <- nil
		return
	}

	response, err := getWithRetry(metaImage.url)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		errch <- err
		done <- false
		return
	}

	output, err := os.Create(metaImage.path)
	if output != nil {
		defer output.Close()
	}
	if err != nil {
		errch <- err
		done <- false
		return
	}
	_, err = io.Copy(output, response.Body)
	if err != nil {
		errch <- err
		done <- false
		return
	}
	errch <- nil
	done <- true
	return
}

func getWithRetry(url string) (*http.Response, error) {
	var (
		err error
		response *http.Response
		retries int = 3
	)
	for retries > 0 {
		response, err = http.Get(url)
		if err != nil {
			retries -= 1
			time.Sleep(5 * time.Millisecond)
		} else {
			break
		}
	}
	return response, err
}

func getProgressBar(max int64, title string) *progressbar.ProgressBar {
	theme := progressbar.Theme{
		Saucer: "=",
		SaucerHead: ">",
		SaucerPadding: " ",
		BarStart: "[",
		BarEnd: "]",
	}
	return progressbar.NewOptions64(
		max,
		progressbar.OptionShowIts(),
		progressbar.OptionSetItsString("photos"),
		progressbar.OptionSetDescription(title),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionShowCount(),
		progressbar.OptionSetTheme(theme),
	)

}
