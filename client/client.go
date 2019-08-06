package client

import (
	"fmt"
	"io/ioutil"
	"loader/dto"
	"loader/utils"
	"net/http"
	"strconv"
	"sync"
)

// TODO: select automatically
const chunkSize = 100
const threadCount = 2

type loader struct {
	client *http.Client
	url    string
	wg     *sync.WaitGroup
	data   chan<- *dto.Chunk
}

func (ld *loader) startThread(offset int, end int) {
	chunkStart := offset
	chunkEnd := chunkStart + chunkSize

	for chunkStart <= end {
		_start := chunkStart
		_end := utils.Min(chunkEnd, end)

		request, err := http.NewRequest("GET", ld.url, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		request.Header.Add("Range", fmt.Sprintf("%d-%d", _start, _end))

		resp, err := ld.client.Do(request)
		if err != nil {
			fmt.Println(err)
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}

		chunk := dto.Chunk{
			Data:  body,
			Start: int64(_start),
		}
		ld.data <- &chunk

		resp.Body.Close()

		chunkStart = chunkEnd
		chunkEnd = chunkStart + chunkSize
	}
	ld.wg.Done()
}

func getSize(client *http.Client, url string) (int, error) {
	request, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	contentLength := resp.Header["Content-Length"]
	var size int
	if len(contentLength) > 0 {
		size, err = strconv.Atoi(resp.Header["Content-Length"][0])
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
	}

	return size, nil
}

// Start - start loading data
func Start(url string, data chan<- *dto.Chunk) int {
	client := &http.Client{}
	wg := sync.WaitGroup{}

	size, err := getSize(client, url)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	go func() {
		fmt.Println("Loading from:", url)
		ld := &loader{
			client: client,
			url:    url,
			wg:     &wg,
			data:   data,
		}

		threadChunkSize := size/threadCount + 1
		offset := 0
		for offset < size {
			end := utils.Min(offset+threadChunkSize, size)

			wg.Add(1)
			go ld.startThread(offset, end)
			offset += threadChunkSize
		}

		wg.Wait()
		close(data)
	}()

	return size
}
