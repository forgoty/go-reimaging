package downloadcommand

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type photoDownloader struct {
	client           *http.Client
	semaphoreChannel chan struct{}
	errorChannel     chan error
}

func NewPhotoDownloader(client *http.Client, semChannel chan struct{}, errCh chan error) *photoDownloader {
	return &photoDownloader{
		client:           client,
		semaphoreChannel: semChannel,
		errorChannel:     errCh,
	}
}

func (d *photoDownloader) Download(url, filename, path string) {
	d.semaphoreAcquire()
	defer d.semaphoreRelease()

	if isFileExists(path) {
		fmt.Println()
		msg := fmt.Sprintf("File exists: %s", path)
		fmt.Println(msg)
		d.writeError(nil)
		return
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		d.writeError(err)
		return
	}
	res, err := d.getWithRetry(req)
	if err != nil {
		d.writeError(err)
		return
	}
	defer res.Body.Close()

	output, err := os.Create(path)
	if err != nil {
		d.writeError(err)
		return

	}
	defer output.Close()
	_, err = io.Copy(output, res.Body)
	if err != nil {
		d.writeError(err)
		return
	}
	d.writeError(nil)
}

func (d *photoDownloader) semaphoreAcquire() {
	d.semaphoreChannel <- struct{}{}
}

func (d *photoDownloader) semaphoreRelease() {
	<-d.semaphoreChannel
}

func (d *photoDownloader) getWithRetry(req *http.Request) (*http.Response, error) {
	var (
		err      error = fmt.Errorf("Cannot download file for url: %s", req.URL)
		response *http.Response
		retries  int = 3
	)
	for retries > 0 {
		response, err = d.client.Do(req)
		if err == nil && response.StatusCode == http.StatusOK {
			break
		}
		retries -= 1
		time.Sleep(5 * time.Millisecond)
	}
	return response, err
}

func (d *photoDownloader) writeError(err error) {
	d.errorChannel <- err
}

func isFileExists(path string) bool {
	_, err := os.Stat(path)
	return os.IsExist(err)
}
