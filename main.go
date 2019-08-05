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

	writer, err := writer.NewFileWriter(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer writer.Close()

	client := client.NewClient(url)
	data, err := client.Start()
	if err != nil {
		fmt.Println(err)
		return
	}

	writer.Write(data)
}
