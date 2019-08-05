package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

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

// Start - start loading data
// TODO: use channel
func (hc *HTTPClient) Start() ([]byte, error) {
	fmt.Println("loading from:", hc.url)
	request, err := http.NewRequest("GET", hc.url, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

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

	return body, nil
}
