package client

import (
	"fmt"
	"io/ioutil"
	"loader/dto"
	"loader/utils"
	"net/http"
	"strconv"
)

// TODO: select automatically
const maxChunkSize = 300
const threadCount = 2

type loader struct {
	client *http.Client
	url    string
	data   chan<- *dto.Chunk
}

func (ld *loader) startThread(descriptor dto.ChunkDescriptor) {
	offset := descriptor.Offset
	cursor := descriptor.Start + offset

	for cursor <= descriptor.End {
		chunkEnd := utils.Min(cursor+maxChunkSize, descriptor.End)
		request, err := http.NewRequest("GET", ld.url, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		request.Header.Add("Range", fmt.Sprintf("%d-%d", cursor, chunkEnd))

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

		chunk := dto.Chunk{ChunkDescriptor: descriptor, Data: body}
		ld.data <- &chunk

		resp.Body.Close()

		cursor = chunkEnd
		descriptor.Offset = cursor - descriptor.Start
	}
}

func (ld *loader) getSize() (int, error) {
	request, err := http.NewRequest("HEAD", ld.url, nil)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	resp, err := ld.client.Do(request)
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

func createChunkDescriptors(size int64) []dto.ChunkDescriptor {
	threadChunkSize := size/threadCount + 1
	var start int64
	chunkDescriptors := make([]dto.ChunkDescriptor, 0)
	for start < size {
		end := utils.Min(start+threadChunkSize, size)
		chunkDescriptor := dto.ChunkDescriptor{
			Start:  int64(start),
			Offset: 0,
			End:    int64(end),
		}
		chunkDescriptors = append(chunkDescriptors, chunkDescriptor)
		start += threadChunkSize
	}

	return chunkDescriptors
}

// Start - start loading data
func Start(info *dto.ProcessDescriptor, data chan<- *dto.Chunk) error {
	ld := &loader{
		client: &http.Client{},
		url:    info.URL,
		data:   data,
	}

	_size, err := ld.getSize()
	if err != nil {
		return err
	}
	size := int64(_size)
	info.Size = size

	chunkDescriptors := createChunkDescriptors(size)

	fmt.Println("Loading from:", info.URL, size, "bytes")

	for _, chunkDescriptor := range chunkDescriptors {
		go ld.startThread(chunkDescriptor)
	}

	return nil
}
