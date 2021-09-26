package uploadcommand

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	vkw "github.com/forgoty/go-reimaging/pkg/reimaging/vkwrapper"
)

type FileGroupUploader struct {
	client           *http.Client
	vkWrapper        vkw.VKWrapper
	semaphoreChannel chan struct{}
	errorChannel     chan error
}

func NewFileGroupUploader(client *http.Client, VKWrapper vkw.VKWrapper, semChannel chan struct{}, errCh chan error) *FileGroupUploader {
	return &FileGroupUploader{
		client:           client,
		vkWrapper:        VKWrapper,
		semaphoreChannel: semChannel,
		errorChannel:     errCh,
	}
}

func (f *FileGroupUploader) Upload(group []string, uploadServer string, albumId int) {
	f.semaphoreAcquire()
	defer f.semaphoreRelease()

	res, err := f.uploadGroup(group, uploadServer)
	if err != nil {
		f.writeError(err, group)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		f.writeError(fmt.Errorf("bad status: %s", res.Status), group)
		return
	}

	body, _ := io.ReadAll(res.Body)
	err = f.vkWrapper.PhotosSave(body, albumId)
	f.writeError(err, group)
}

func (f *FileGroupUploader) semaphoreAcquire() {
	f.semaphoreChannel <- struct{}{}
}

func (f *FileGroupUploader) semaphoreRelease() {
	<-f.semaphoreChannel
}

func (f *FileGroupUploader) uploadGroup(group []string, uploadServer string) (*http.Response, error) {
	body := bytes.Buffer{}
	writer := multipart.NewWriter(&body)

	err := prepareMultipartBody(group, writer, &body)
	if err != nil {
		return nil, err
	}
	writer.Close()

	req, err := http.NewRequest("POST", uploadServer, &body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return f.client.Do(req)
}

func (f *FileGroupUploader) writeError(err error, group []string) {
	if err != nil {
		fmt.Println()
		fmt.Printf("File Upload Failed for group:")
		fmt.Println()
		for _, path := range group {
			fmt.Println(path)
		}
		fmt.Println(err)
		fmt.Println()
	}
	f.errorChannel <- err
}

func prepareMultipartBody(group []string, writer *multipart.Writer, bf *bytes.Buffer) error {
	for i := range group {
		r, err := os.ReadFile(group[i])
		if err != nil {
			return err
		}
		w, err := writer.CreateFormFile(fmt.Sprintf("file%d", i+1), group[i])
		if err != nil {
			return err
		}
		if _, err = io.Copy(w, bytes.NewReader(r)); err != nil {
			return err
		}
	}
	return nil
}
