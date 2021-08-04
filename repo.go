package acc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

const (
	_defaultContentType = "application/json"
	_errContext         = "http_repo"
)

type (
	payload struct {
		Data *account `json:"data"`
	}

	marshaler func(v interface{}) ([]byte, error)

	repo struct {
		addr string
		port string
		marshaler
	}
)

func newRepo(addr, port string) repo {
	return repo{
		addr:      addr,
		port:      port,
		marshaler: json.Marshal,
	}
}

func (r repo) create(acc account) (*account, error) {
	const success = 201
	wrapErr := func(err error, msg string) error {
		return errors.Wrapf(err, "%s create %s", _errContext, msg)
	}

	data, err := r.marshaler(payload{Data: &acc})
	if err != nil {
		return nil, wrapErr(err, "marshal")
	}

	url := fmt.Sprintf("http://%s:%s/v1/organisation/accounts", r.addr, r.port)
	// TODO: use mock for it HARD
	resp, err := http.Post(url, _defaultContentType, bytes.NewBuffer(data))
	if err != nil {
		return nil, wrapErr(err, "request")
	}

	// TODO: use mock for it HARD
	if resp.StatusCode != success {
		return nil, wrapErr(err, "not success != 201")
	}

	var ret payload
	// TODO: use mock for it
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&ret); err != nil {
		return nil, wrapErr(err, "decode")
	}

	return ret.Data, nil
}

// TODO: improve it
func (r repo) health() error {
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
