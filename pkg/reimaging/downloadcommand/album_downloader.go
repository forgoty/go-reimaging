package downloadcommand

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/forgoty/go-reimaging/pkg/reimaging/progressbar"
	vkw "github.com/forgoty/go-reimaging/pkg/reimaging/vkwrapper"
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
	urlCalculator := NewUrlCalculator(ad.vkWrapper)
	urls := urlCalculator.Calculate(album)
	if len(urls) == 0 {
		return
	}
	pathDir := createAlbumDir(ad.options.DownloadPath, album.Title)
	downloadPhotos(urls, pathDir, album.Title)
}

func createAlbumDir(path, title string) string {
	pathDir := filepath.Join(path, title)
	if _, err := os.Stat(pathDir); os.IsNotExist(err) {
		if err = os.MkdirAll(pathDir, os.ModePerm); err != nil {
			fmt.Printf("Cannot create a folder %s: %s", pathDir, err)
			fmt.Println()
			os.Exit(1)
		}
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
				err = bar.Finish()
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println()
				fmt.Println("Too many errors occured recently")
				os.Exit(1)
			}
		}
		err = bar.Add(1)
		if err != nil {
			fmt.Println(err)
		}
		results = append(results, err)
		if len(results) == total {
			err = bar.Finish()
			if err != nil {
				fmt.Println(err)
			}
			break
		}
	}
}
