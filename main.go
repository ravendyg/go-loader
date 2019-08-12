package main

import (
	"flag"
	"fmt"
	"loader/client"
	"loader/dto"
	"loader/writer"
	"strings"
)

func main() {
	ftprURL := flag.String("url", "", "Path to the file to be downloaded")
	ftprName := flag.String("name", "", "File name")
	flag.Parse()
	url := *ftprURL
	fileName := *ftprName
	if url == "" {
		url = "http://localhost:3011/"
	}

	dataChannel := make(chan dto.Chunk)

	if fileName == "" {
		nameChunks := strings.Split(url, "/")
		fileName = nameChunks[len(nameChunks)-1]
	}
	if fileName == "" {
		fileName = "download"
	}

	fi, err := writer.NewFileWriter(url, fileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	descriptor, err := fi.ReadMeatada()
	if err != nil {
		fmt.Println(err)
		return
	}

	if descriptor == nil {
		descriptor = &dto.ProcessDescriptor{
			URL:              url,
			FileName:         fileName,
			ChunkDescriptors: nil,
			Size:             0,
		}
	}

	size, chunkDescriptors, err := client.Start(descriptor, dataChannel)
	if err != nil {
		fmt.Println(err)
		return
	}

	descriptor.Size = size
	descriptor.ChunkDescriptors = chunkDescriptors

	for chunk := range dataChannel {
		fi.WriteData(&chunk)
		var loaded int64
		// TODO: use map?
		for i := 0; i < len(descriptor.ChunkDescriptors); i++ {
			if chunkDescriptors[i].ID == chunk.ID {
				chunkDescriptors[i].Offset = chunk.Cursor - chunkDescriptors[i].Start + int64(len(chunk.Data))
			}

			loaded += chunkDescriptors[i].Offset
		}
		fi.WriteMetaData(descriptor)
		fmt.Printf("Loaded %d%%\n", loaded*100/descriptor.Size)
		if loaded >= descriptor.Size {
			close(dataChannel)
		}
	}

	err = fi.Finish()
	if err != nil {
		fmt.Println(err)
	}
}
