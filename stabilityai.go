package stabilityai

import "net/http"

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	httpClient HTTPClient
	apiKey     string
}

func NewClient(httpClient HTTPClient, apiKey string) *Client {
	return &Client{
		httpClient: httpClient,
		apiKey:     apiKey,
	}
}
