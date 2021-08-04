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
	_defaultContentType = "application/json"
	_errContext         = "http_repo"
)

type (
	payload struct {
		Data *account `json:"data"`
	}

	marshal func(v interface{}) ([]byte, error)
	post    func(url, contentType string, body io.Reader) (resp *http.Response, err error)
	decode  func(d *json.Decoder, v interface{}) error

	httpRepository struct {
		addr string
		port string
		marshal
		post
		decode
	}
)

// TODO:
// - make it exported;
// - Add documentation;
// - Change parameters to Functional Options: https://github.com/uber-go/guide/blob/master/style.md#functional-options;
func newHTTPRepository(addr, port string) httpRepository {
	return httpRepository{
		addr:    addr,
		port:    port,
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
