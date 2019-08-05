package main

import (
	"flag"
	"fmt"
	"loader/client"
	"loader/writer"
)

func main() {
	ftpr := flag.String("url", "", "Path to the file to be downloaded")
	flag.Parse()
	url := *ftpr

	dataChannel := make(chan *client.Chunk)

	writer, err := writer.NewFileWriter(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer writer.Close()

	go client.Start(url, dataChannel)

	for chunk := range dataChannel {
		writer.Write(chunk.Data, chunk.Start)
	}
}
