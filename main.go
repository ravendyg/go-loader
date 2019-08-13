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

	loader := client.NewClient(url, dataChannel)
	size, err := loader.GetSize()
	if err != nil {
		fmt.Println(err)
		return
	}

	if descriptor == nil || descriptor.Size != size {
		descriptor = &dto.ProcessDescriptor{
			URL:              url,
			FileName:         fileName,
			ChunkDescriptors: loader.CreateChunkDescriptors(size),
			Size:             size,
		}
	}

	loader.Start(descriptor)

	for chunk := range dataChannel {
		fi.WriteData(&chunk)
		var loaded int64
		for i := 0; i < len(descriptor.ChunkDescriptors); i++ {
			if descriptor.ChunkDescriptors[i].ID == chunk.ID {
				descriptor.ChunkDescriptors[i].Offset = chunk.Cursor -
					descriptor.ChunkDescriptors[i].Start + chunk.Size
			}

			loaded += descriptor.ChunkDescriptors[i].Offset
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
