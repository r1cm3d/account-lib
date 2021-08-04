package acc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

const (
	_defaultHTTPAddress = "0.0.0.0"
	_defaultHTTPPort    = "8080"
	_defaultContentType = "application/json"
	_errContext         = "http_repo"
)

type (
	httpOption interface {
		apply(*httpOptions)
	}
	httpOptions struct {
		addr string
		port string
	}
	addrOption string
	portOption string

	payload struct {
		Data *account `json:"data"`
	}

	marshal        func(v interface{}) ([]byte, error)
	post           func(url, contentType string, body io.Reader) (resp *http.Response, err error)
	decode         func(d *json.Decoder, v interface{}) error
	httpRepository struct {
		addr string
		port string
		marshal
		post
		decode
	}
)

// WithAddr attaches server address to HTTP client.
// Default is 0.0.0.0
func WithAddr(addr string) httpOption {
	return addrOption(addr)
}

// WithPort attaches server TCP port to HTTP client.
// Default is 8080
func WithPort(port string) httpOption {
	return portOption(port)
}

// NewHTTPRepository instantiate a httpRepository based on httpOption(s) passed as arguments.
//
// Example:
//  repository := NewHTTPRepository(acc.WithPort("8080"))
func NewHTTPRepository(opts ...httpOption) httpRepository {
	options := httpOptions{
		addr: _defaultHTTPAddress,
		port: _defaultHTTPPort,
	}

	for _, o := range opts {
		o.apply(&options)
	}

	return httpRepository{
		addr:    options.addr,
		port:    options.port,
		marshal: json.Marshal,
		post:    http.Post,
		decode:  func(d *json.Decoder, v interface{}) error { return d.Decode(v) },
	}
}

func (r httpRepository) create(acc account) (*account, error) {
	const (
		success = 201
		urlBase = "http://%s:%s/v1/organisation/accounts"
	)
	wrapErr := func(err error, msg string) error {
		return errors.Wrapf(err, "%s create %s", _errContext, msg)
	}

	data, err := r.marshal(payload{Data: &acc})
	if err != nil {
		return nil, wrapErr(err, "marshal")
	}

	url := fmt.Sprintf(urlBase, r.addr, r.port)
	resp, err := r.post(url, _defaultContentType, bytes.NewBuffer(data))
	if err != nil {
		return nil, wrapErr(err, "request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != success {
		return nil, wrapErr(errors.New("not success != 201"), "status code verification")
	}

	var ret payload
	dec := json.NewDecoder(resp.Body)
	if err := r.decode(dec, &ret); err != nil {
		return nil, wrapErr(err, "decode")
	}

	return ret.Data, nil
}

// TODO: improve it
func (r httpRepository) health() error {
	// TODO: use mock for it
	resp, err := http.Get(fmt.Sprintf("http://%s:%s/v1/health", r.addr, r.addr))

	if err != nil {
		return err
	}
	// TODO: use mock for it
	defer resp.Body.Close()
	var data map[string]interface{}
	// TODO: use mock for it
	json.NewDecoder(resp.Body).Decode(&data)

	// TODO: unmarshal status and try status code
	fmt.Println(data)
	return nil
}

func (a addrOption) apply(opts *httpOptions) {
	opts.addr = string(a)
}
func (p portOption) apply(opts *httpOptions) {
	opts.port = string(p)
}
