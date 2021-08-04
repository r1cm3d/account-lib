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
)

type (
	marshal        func(v interface{}) ([]byte, error)
	post           func(url, contentType string, body io.Reader) (resp *http.Response, err error)
	decode         func(d *json.Decoder, v interface{}) error
	httpRepository struct {
		marshal
		post
		decode

		addr     string
		port     string
		errCtx   string
		contType string
	}
)

type (
	payload struct {
		Data *data `json:"data"`
	}
	data struct {
		Attributes     *attributes `json:"attributes,omitempty"`
		ID             string      `json:"id,omitempty"`
		OrganisationID string      `json:"organisation_id,omitempty"`
		Type           string      `json:"type,omitempty"`
		Version        *int64      `json:"version,omitempty"`
	}
	attributes struct {
		Classification          *string  `json:"account_classification,omitempty"`
		MatchingOptOut          *bool    `json:"account_matching_opt_out,omitempty"`
		Number                  string   `json:"account_number,omitempty"`
		AlternativeNames        []string `json:"alternative_names,omitempty"`
		BankID                  string   `json:"bank_id,omitempty"`
		BankIDCode              string   `json:"bank_id_code,omitempty"`
		BaseCurrency            string   `json:"base_currency,omitempty"`
		Bic                     string   `json:"bic,omitempty"`
		Country                 *string  `json:"country,omitempty"`
		Iban                    string   `json:"iban,omitempty"`
		JointAccount            *bool    `json:"joint_account,omitempty"`
		Name                    []string `json:"name,omitempty"`
		SecondaryIdentification string   `json:"secondary_identification,omitempty"`
		Status                  *string  `json:"status,omitempty"`
		Switched                *bool    `json:"switched,omitempty"`
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
		addr:     options.addr,
		port:     options.port,
		errCtx:   "http_repository",
		contType: "application/json",
		marshal:  json.Marshal,
		post:     http.Post,
		decode:   func(d *json.Decoder, v interface{}) error { return d.Decode(v) },
	}
}

func (r httpRepository) create(acc data) (*data, error) {
	const (
		success = 201
		urlBase = "http://%s:%s/v1/organisation/accounts"
	)
	wrapErr := func(err error, msg string) error {
		return errors.Wrapf(err, "%s create %s", r.errCtx, msg)
	}

	data, err := r.marshal(payload{Data: &acc})
	if err != nil {
		return nil, wrapErr(err, "marshal")
	}

	url := fmt.Sprintf(urlBase, r.addr, r.port)
	resp, err := r.post(url, r.contType, bytes.NewBuffer(data))
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
