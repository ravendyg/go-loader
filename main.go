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

	size := client.Start(url, dataChannel)

	loaded := 0
	for chunk := range dataChannel {
		writer.WriteData(chunk)
		loaded += len(chunk.Data)
		fmt.Printf("Loaded %d%%\n", loaded*100/size)
	}
	err = writer.Finish()
	if err != nil {
		fmt.Println(err)
	}
}
