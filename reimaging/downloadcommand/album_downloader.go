package downloadcommand

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/forgoty/go-reimaging/reimaging/progressbar"
	vkw "github.com/forgoty/go-reimaging/reimaging/vkwrapper"
)

type downloadOptions struct {
	UserID       int
	DownloadPath string
	NeedSystem   bool
}

func NewDownloadOptions(userID int, downloadPath string, needSystem bool) *downloadOptions {
	return &downloadOptions{
		UserID:       userID,
		DownloadPath: downloadPath,
		NeedSystem:   needSystem,
	}
}

type AlbumDownloader struct {
	vkWrapper vkw.VKWrapper
	options   *downloadOptions
}

func NewAlbumDownloader(vk vkw.VKWrapper, options *downloadOptions) *AlbumDownloader {
	return &AlbumDownloader{
		vkWrapper: vk,
		options:   options,
	}
}

func (ad *AlbumDownloader) DownloadAlbumByID(albumID int) {
	albums := ad.vkWrapper.GetAlbums(ad.options.UserID, ad.options.NeedSystem)
	for _, album := range albums {
		if album.ID == albumID {
			ad.DownloadAlbum(album)
		}
	}
}

func (ad *AlbumDownloader) DownloadAll() {
	albums := ad.vkWrapper.GetAlbums(ad.options.UserID, ad.options.NeedSystem)
	for _, album := range albums {
		ad.DownloadAlbum(album)
	}
}

func (ad *AlbumDownloader) DownloadAlbum(album vkw.PhotoAlbum) {
	if album.Size > 1000 {
		fmt.Println("Calculating urls...")
	}
	offsets := getOffset(album.Size)
	photosUrls := []string{}
	for _, offset := range offsets {
		photosUrls = append(photosUrls, ad.vkWrapper.GetPhotoURLs(album, offset)...)
	}
	if len(photosUrls) > 0 {
		pathDir := createAlbumDir(ad.options.DownloadPath, album.Title)
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
	client := http.Client{}
	semaphoreChan := make(chan struct{}, 200)
	errCh := make(chan error)
	defer func() {
		close(semaphoreChan)
		close(errCh)
	}()

	for _, url := range photosUrls {
		metaImage := getPhotoMetaData(url, pathDir)
		photoDownloader := NewPhotoDownloader(&client, semaphoreChan, errCh)
		go photoDownloader.Download(metaImage.url, metaImage.filename, metaImage.path)
	}

	bar := progressbar.NewProgressBar(photosCount, albumTitle)
	waitDownload(errCh, bar, photosCount)
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

func getFileName(url string) string {
	split := strings.Split(url, "/")
	filename := strings.ReplaceAll(split[len(split)-1], "-", "")
	if index := strings.Index(filename, "?"); index > 0 {
		filename = filename[0:index]
	}
	return filename
}

func waitDownload(errCh chan error, bar progressbar.ProgressBarHandler, total int) {
	var results []error
	var errors []error
	for {
		err := <-errCh
		if err != nil {
			errors = append(errors, err)
			if len(errors) > 3 {
				bar.Finish()
				fmt.Println()
				fmt.Println("Too many errors occured recently")
				os.Exit(1)
			}
		}
		bar.Add(1)
		results = append(results, err)
		if len(results) == total {
			bar.Finish()
			break
		}
	}
}
