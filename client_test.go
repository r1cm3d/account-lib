package account

import (
	"github.com/pkg/errors"
	"io"
	"net/http"
	"testing"
)

type mockRequesterError struct{}

var (
	_clientWithBuildRequestError = httpclient{
		errCtx: "http_client",
		buildRequest: func(method, url string, body io.Reader) (*http.Request, error) {
			return nil, errors.New("error on buildRequest")
		},
	}
	_clientWithDoError = httpclient{
		errCtx:       "http_client",
		buildRequest: func(method, url string, body io.Reader) (*http.Request, error) { return nil, nil },
		requester:    mockRequesterError{},
	}
)

func TestRequest_Error(t *testing.T) {
	cases := []struct {
		name string
		in   httpclient
		want error
	}{
		{"buildRequest", _clientWithBuildRequestError, errors.New("http_client#request() buildRequest: error on buildRequest")},
		{"do", _clientWithDoError, errors.New("http_client#request() do: error on do")},
	}
	for _, tt := range cases {
		_, got := tt.in.request("method", "url", nil)
		if got.Error() != tt.want.Error() {
			t.Errorf("Request_Error(%v) got: %v, want: %v", tt.name, got, tt.want)
		}
	}
}

func (m mockRequesterError) Do(_ *http.Request) (*http.Response, error) {
	return nil, errors.New("error on do")
}
