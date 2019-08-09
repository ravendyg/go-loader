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

	writer, err := writer.NewFileWriter(url, fileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	descriptor := dto.ProcessDescriptor{
		URL:      url,
		FileName: fileName,
	}

	client.Start(&descriptor, dataChannel)

	var loaded int64
	for chunk := range dataChannel {
		writer.WriteData(chunk)
		loaded += int64(len(chunk.Data))
		fmt.Printf("Loaded %d%%\n", loaded*100/descriptor.Size)
		if loaded >= descriptor.Size {
			close(dataChannel)
		}
	}

	err = writer.Finish()
	if err != nil {
		fmt.Println(err)
	}
}
