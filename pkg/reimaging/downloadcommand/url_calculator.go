package downloadcommand

import (
	"fmt"
	vkw "github.com/forgoty/go-reimaging/pkg/reimaging/vkwrapper"
	"math"
	"sync"
)

type vkUrlCalculator struct {
	vkWrapper vkw.VKWrapper
}

func NewUrlCalculator(vkWrapper vkw.VKWrapper) *vkUrlCalculator {
	return &vkUrlCalculator{
		vkWrapper: vkWrapper,
	}
}

func (c *vkUrlCalculator) Calculate(album vkw.PhotoAlbum) []string {
	offsets := getOffset(album.Size)
	lenOffsets := len(offsets)
	photosUrls := []string{}
	var wg sync.WaitGroup
	wg.Add(lenOffsets)
	queue := make(chan []string, lenOffsets)

	if album.Size > 1000 {
		fmt.Println("Calculating urls...")
	}
	for _, offset := range offsets {
		go func(offsetNumber int) {
			queue <- c.vkWrapper.GetPhotoURLs(album, offsetNumber)
		}(offset)
	}
	go func() {
		for urls := range queue {
			photosUrls = append(photosUrls, urls...)
			wg.Done()
		}
	}()
	wg.Wait()
	return photosUrls
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
