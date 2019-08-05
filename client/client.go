package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
)

const chunkSize = 20

// Chunk - DTO
type Chunk struct {
	Start int64
	Data  []byte
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
func Start(url string, data chan<- *Chunk) {
	client := &http.Client{}
	wg := sync.WaitGroup{}

	size, err := getSize(client, url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("File size: ", size)
	fmt.Println("Loading from:", url)

	var start int
	var end = start + chunkSize

	for start <= size {
		_start := start
		_end := end
		wg.Add(1)
		go func() {
			request, err := http.NewRequest("GET", url, nil)
			if err != nil {
				fmt.Println(err)
				return
			}

			request.Header.Add("Range", fmt.Sprintf("%d-%d", _start, _end))

			resp, err := client.Do(request)
			if err != nil {
				fmt.Println(err)
				return
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
				return
			}

			chunk := Chunk{
				Data:  body,
				Start: int64(_start),
			}
			data <- &chunk

			resp.Body.Close()
			wg.Done()
		}()
		start = end
		end = start + chunkSize
	}

	wg.Wait()
	close(data)
}
