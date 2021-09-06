package vkwrapper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
)

type VKWrapper interface {
	GetPhotoURLs(album PhotoAlbum, offset int) []string
	GetAlbums(userID int, NeedSystem bool) []PhotoAlbum
	GetAlbumsByAlbumIds(albumIDs []int) []PhotoAlbum
	CreateAlbum(title string) PhotoAlbum
	GetUploadServer(id int) string
	UploadPhoto(albumID int, file io.Reader) error
	UploadFileGroup(client *http.Client, group []string, uploadServer string, albumId int, semCh chan struct{}, errCh chan error)
}

type VkAPIWrapper struct {
	vk *api.VK
}

func NewVKWrapper() *VkAPIWrapper {
	return &VkAPIWrapper{getVk()}
}

func (vkw *VkAPIWrapper) GetPhotoURLs(album PhotoAlbum, offset int) []string {
	responseParams := params.NewPhotosGetBuilder()
	responseParams.OwnerID(album.OwnerID)
	responseParams.AlbumID(strconv.Itoa(album.ID))
	responseParams.Offset(offset)
	responseParams.PhotoSizes(true)
	responseParams.Count(1000)
	response, vkErr := vkw.vk.PhotosGet(responseParams.Params)
	if vkErr != nil {
		fmt.Println(vkErr)
	}

	urls := []string{}
	for _, photo := range response.Items {
		if len(photo.Sizes) > 0 {
			// the most high-rez picture is always last
			urls = append(urls, photo.Sizes[len(photo.Sizes)-1].BaseImage.URL)
		}
	}
	return urls
}

func (vkw *VkAPIWrapper) GetAlbums(userID int, NeedSystem bool) []PhotoAlbum {
	responseParams := params.NewPhotosGetAlbumsBuilder()
	responseParams.OwnerID(userID)
	responseParams.NeedSystem(NeedSystem)
	response, vkErr := vkw.vk.PhotosGetAlbums(responseParams.Params)
	if vkErr != nil {
		fmt.Println(vkErr)
		os.Exit(1)
	}
	albums := []PhotoAlbum{}
	for _, rawAlbum := range response.Items {
		albums = append(
			albums,
			PhotoAlbum{
				ID:      rawAlbum.ID,
				OwnerID: rawAlbum.OwnerID,
				Size:    rawAlbum.Size,
				Title:   rawAlbum.Title,
			},
		)
	}

	return albums
}

func (vkw *VkAPIWrapper) GetAlbumsByAlbumIds(albumIDs []int) []PhotoAlbum {
	responseParams := params.NewPhotosGetAlbumsBuilder()
	responseParams.AlbumIDs(albumIDs)
	response, vkErr := vkw.vk.PhotosGetAlbums(responseParams.Params)
	if vkErr != nil {
		fmt.Println(vkErr)
		os.Exit(1)
	}
	albums := []PhotoAlbum{}
	for _, rawAlbum := range response.Items {
		albums = append(
			albums,
			PhotoAlbum{
				ID:      rawAlbum.ID,
				OwnerID: rawAlbum.OwnerID,
				Size:    rawAlbum.Size,
				Title:   rawAlbum.Title,
			},
		)
	}

	return albums
}

func (vkw *VkAPIWrapper) CreateAlbum(title string) PhotoAlbum {
	responseParams := params.NewPhotosCreateAlbumBuilder()
	responseParams.Title(title)
	responseParams.PrivacyView([]string{"only_me"})
	rawAlbum, vkErr := vkw.vk.PhotosCreateAlbum(responseParams.Params)
	if vkErr != nil {
		fmt.Println(vkErr)
		os.Exit(1)
	}

	return PhotoAlbum{
		ID:      rawAlbum.ID,
		OwnerID: rawAlbum.OwnerID,
		Size:    rawAlbum.Size,
		Title:   rawAlbum.Title,
	}
}

func (vkw *VkAPIWrapper) GetUploadServer(id int) string {
	responseParams := params.NewPhotosGetUploadServerBuilder()
	responseParams.AlbumID(id)
	uploadServerResponse, vkErr := vkw.vk.PhotosGetUploadServer(responseParams.Params)
	if vkErr != nil {
		fmt.Println(vkErr)
		os.Exit(1)
	}
	return uploadServerResponse.UploadURL
}

// Refactor me
func (vkw *VkAPIWrapper) UploadFileGroup(client *http.Client, group []string, uploadServer string, albumId int, semCh chan struct{}, errCh chan error) {
	semCh <- struct{}{}
	b := bytes.Buffer{}
	writer := multipart.NewWriter(&b)
	for i := range group {
		r, err := os.ReadFile(group[i])
		if err != nil {
			errCh <- err
			<-semCh
			return
		}
		w, err := writer.CreateFormFile(fmt.Sprintf("file%d", i+1), group[i])
		if err != nil {
			errCh <- err
			<-semCh
			return
		}
		if _, err = io.Copy(w, bytes.NewReader(r)); err != nil {
			errCh <- err
			<-semCh
			return
		}
	}
	writer.Close()
	req, err := http.NewRequest("POST", uploadServer, &b)
	if err != nil {
		errCh <- err
		<-semCh
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := client.Do(req)
	if err != nil {
		errCh <- err
		<-semCh
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
		errCh <- err
		<-semCh
		return
	}

	body, _ := io.ReadAll(res.Body)
	var handler object.PhotosPhotoUploadResponse
	err = json.Unmarshal(body, &handler)
	if err != nil {
		errCh <- err
		<-semCh
		return
	}

	_, err = vkw.vk.PhotosSave(api.Params{
		"server":      handler.Server,
		"photos_list": handler.PhotosList,
		"aid":         handler.AID,
		"hash":        handler.Hash,
		"album_id":    albumId,
	})
	<-semCh
	errCh <- err
}

func (vkw *VkAPIWrapper) UploadPhoto(albumID int, file io.Reader) error {
	_, err := vkw.vk.UploadPhoto(albumID, file)
	return err
}

func getVk() *api.VK {
	return api.NewVK(os.Getenv("VK_TOKEN"))
}

type PhotoAlbum struct {
	ID      int
	OwnerID int
	Size    int
	Title   string
}
