package downloadcommand

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	progressbar "github.com/schollz/progressbar/v3"
	vkw "github.com/forgoty/go-reimaging/reimaging/vkwrapper"
)

type AlbumDownloader struct {
	VKWrapper *vkw.VKWrapper
	DownloadPath string
	NeedSystem bool
}

func NewAlbumDownloader(userID int, downloadPath string, needSystem bool) *AlbumDownloader {
	return &AlbumDownloader{
		VKWrapper: vkw.NewVKWrapper(userID),
		DownloadPath: downloadPath,
		NeedSystem: needSystem,
	}
}

func (ad *AlbumDownloader) DownloadAlbumByID(albumID int) {
	albums := ad.VKWrapper.GetAlbums(ad.NeedSystem)
	for _, album := range albums {
		if album.ID == albumID {
			ad.DownloadAlbum(album)
		}
	}
}

func (ad *AlbumDownloader) DownloadAll() {
	albums := ad.VKWrapper.GetAlbums(ad.NeedSystem)
	for _, album := range albums {
		ad.DownloadAlbum(album)
	}
}

func (ad *AlbumDownloader) DownloadAlbum(album vkw.PhotoAlbum) {
	offsets := getOffset(album.Size)
	photosUrls := []string{}
	for _, offset := range offsets {
		photosUrls = append(photosUrls, ad.VKWrapper.GetPhotoURLs(album, offset)...)
	}
	if len(photosUrls) > 0 {
		pathDir := createAlbumDir(ad.DownloadPath, album.Title)
		downloadPhotos(photosUrls, pathDir, album.Title)
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

func createAlbumDir(path, title string) string {
	pathDir := filepath.Join(path, title)
	if _, serr := os.Stat(pathDir); serr != nil {
		os.MkdirAll(pathDir, os.ModePerm)
	}
	return pathDir
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

type metaDownloadedPhoto struct {
	url      string
	filename string
	path     string
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
