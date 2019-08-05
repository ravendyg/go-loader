package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

const chunkSize = 20

// HTTPClient - client interface
type HTTPClient struct {
	client *http.Client
	url    string
}

// NewClient - create new client
func NewClient(url string) *HTTPClient {
	client := &http.Client{}
	return &HTTPClient{client, url}
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
// TODO: use channel
func (hc *HTTPClient) Start() ([]byte, error) {
	size, err := getSize(hc.client, hc.url)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println("File size: ", size)
	fmt.Println("loading from:", hc.url)
	buf := make([]byte, 0)

	var start int
	var end = start + chunkSize

	for end <= size {

		request, err := http.NewRequest("GET", hc.url, nil)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		request.Header.Add("Range", fmt.Sprintf("%d-%d", start, end))

		resp, err := hc.client.Do(request)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		buf = append(buf, body...)

		start = end + 1
		end = start + chunkSize
	}

	return buf, nil
}
