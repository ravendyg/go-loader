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

	dataChannel := make(chan *dto.Chunk)

	if fileName == "" {
		nameChunks := strings.Split(url, "/")
		fileName = nameChunks[len(nameChunks)-1]
	}

	fi, err := writer.NewFileWriter(url, fileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	a, err := fi.ReadMeatada()
	fmt.Println(a, err)

	descriptor := dto.ProcessDescriptor{
		URL:              url,
		FileName:         fileName,
		ChunkDescriptors: nil,
		Size:             0,
		Loaded:           0,
	}

	size, chunkDescriptors, err := client.Start(&descriptor, dataChannel)
	if err != nil {
		fmt.Println(err)
		return
	}

	descriptor.Size = size
	descriptor.ChunkDescriptors = chunkDescriptors

	for chunk := range dataChannel {
		fi.WriteData(chunk)
		descriptor.Loaded += int64(len(chunk.Data))
		// TODO: use map?
		for i := 0; i < len(descriptor.ChunkDescriptors); i++ {
			if chunkDescriptors[i].ID == chunk.ChunkDescriptor.ID {
				chunkDescriptors[i] = chunk.ChunkDescriptor
			}
		}
		fmt.Printf("Loaded %d%%\n", descriptor.Loaded*100/descriptor.Size)
		if descriptor.Loaded >= descriptor.Size {
			close(dataChannel)
		}
	}

	err = fi.Finish()
	if err != nil {
		fmt.Println(err)
	}
}
