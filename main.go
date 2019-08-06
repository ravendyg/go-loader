package main

import (
	"flag"
	"fmt"
	"loader/client"
	"loader/dto"
	"loader/writer"
)

func main() {
	ftpr := flag.String("url", "", "Path to the file to be downloaded")
	flag.Parse()
	url := *ftpr

	dataChannel := make(chan *dto.Chunk)

	writer, err := writer.NewFileWriter(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer writer.Close()

	size := client.Start(url, dataChannel)

	loaded := 0
	for chunk := range dataChannel {
		writer.Write(chunk.Data, chunk.Start)
		loaded += len(chunk.Data)
		fmt.Printf("Loaded %d%%\n", loaded*100/size)
	}
}
