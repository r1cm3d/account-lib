package account

import (
	"io"
	"net/http"
)

type (
	method       string
	buildRequest func(method, url string, body io.Reader) (*http.Request, error)
	requester    interface {
		Do(req *http.Request) (*http.Response, error)
	}
	httpclient struct {
		buildRequest
		requester
	}
)

const (
	_post   = "POST"
	_delete = "DELETE"
	_get    = "GET"
)

func newHttpClient() httpclient {
	return httpclient{
		buildRequest: http.NewRequest,
		requester:    &http.Client{},
	}
}

func (h httpclient) request(method method, url string, body io.Reader) (*http.Response, error) {
	req, err := h.buildRequest(string(method), url, body)
	if err != nil {
		// TODO: wrap it
		return nil, err
	}

	resp, err := h.Do(req)
	if err != nil {
		// TODO: wrap it
		return nil, err
	}

	return resp, nil
}
