package account

import (
	"github.com/pkg/errors"
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

		errCtx string
	}
)

const (
	_post   = "POST"
	_delete = "DELETE"
	_get    = "GET"
)

func newHTTPClient() httpclient {
	return httpclient{
		buildRequest: http.NewRequest,
		requester:    &http.Client{},
		errCtx:       "http_client",
	}
}

func (h httpclient) request(method method, url string, body io.Reader) (*http.Response, error) {
	req, err := h.buildRequest(string(method), url, body)
	if err != nil {
		return nil, errors.Wrapf(err, "%s#request() buildRequest", h.errCtx)
	}

	resp, err := h.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "%s#request() do", h.errCtx)
	}

	return resp, nil
}
