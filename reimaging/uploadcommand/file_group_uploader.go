package uploadcommand

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	vkw "github.com/forgoty/go-reimaging/reimaging/vkwrapper"
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
		f.writeOutput(err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		f.writeOutput(fmt.Errorf("bad status: %s", res.Status))
		return
	}

	body, _ := io.ReadAll(res.Body)
	err = f.vkWrapper.PhotosSave(body, albumId)
	f.writeOutput(err)
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
	defer writer.Close()

	err := prepareMultipartBody(group, writer, &body)
	if err != nil {
		return &http.Response{}, err
	}

	req, err := http.NewRequest("POST", uploadServer, &body)
	if err != nil {
		return &http.Response{}, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return f.client.Do(req)
}

func (f *FileGroupUploader) writeOutput(err error) {
	if err != nil {
		fmt.Println()
		fmt.Println(err)
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
