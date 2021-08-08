package account

import (
	"io"
	"net/http"
)

type (
	method     string
	httpclient struct {
	}
)

const (
	_post   = "POST"
	_delete = "DELETE"
	_get    = "GET"
)

func (h httpclient) request(method method, url string, body io.Reader) (*http.Response, error) {
	client := &http.Client{}

	req, err := http.NewRequest(string(method), url, body)
	if err != nil {
		// TODO: wrap it
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		// TODO: wrap it
		return nil, err
	}

	return resp, nil
}
