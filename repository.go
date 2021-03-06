package account

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
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
	client interface {
		request(method method, url string, body io.Reader) (*http.Response, error)
	}
	marshal        func(v interface{}) ([]byte, error)
	decode         func(d *json.Decoder, v interface{}) error
	httpRepository struct {
		marshal
		client
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
		ID             string      `json:"ID,omitempty"`
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

const (
	_defaultHTTPAddress = "0.0.0.0"
	_defaultHTTPPort    = "8080"
)

// NewHTTPRepository instantiates a httpRepository based on httpOption(s) passed as arguments. If no argument is passed
// the defaults will be used.
//
// Example:
//  repository := NewHTTPRepository(acc.WithPort("8080"))
//  repository := NewHTTPRepository(acc.WithPort("8080"), acc.WithAddr("0.0.0.0"))
//  repository := NewHTTPRepository()
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
		client:   newHTTPClient(),
		decode:   func(d *json.Decoder, v interface{}) error { return d.Decode(v) },
	}
}

// WithAddr attaches server address to HTTP client.
// Default is 0.0.0.0
//
// See: https://github.com/uber-go/guide/blob/master/style.md#functional-options
func WithAddr(addr string) httpOption {
	return addrOption(addr)
}

// WithPort attaches server TCP port to HTTP client.
// Default is 8080
//
// See: https://github.com/uber-go/guide/blob/master/style.md#functional-options
func WithPort(port string) httpOption {
	return portOption(port)
}

func (r httpRepository) create(acc data) (*data, error) {
	const urlBase = "http://%s:%s/v1/organisation/accounts"
	wrapErr := func(err error, msg string) error {
		return errors.Wrapf(err, "%s#create() %s", r.errCtx, msg)
	}

	data, err := r.marshal(payload{Data: &acc})
	if err != nil {
		return nil, wrapErr(err, "marshal")
	}

	url := fmt.Sprintf(urlBase, r.addr, r.port)
	resp, err := r.request(_post, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, wrapErr(err, "request")
	}
	defer resp.Body.Close()

	return r.handleCreateResp(resp)
}

func (r httpRepository) handleCreateResp(resp *http.Response) (*data, error) {
	const (
		success     = 201
		clientError = 400
	)

	switch resp.StatusCode {
	case success:
		return r.parseSuccess(resp.Body)
	case clientError:
		return r.parseClientError(resp.Body)
	default:
		return nil, errors.New(fmt.Sprintf("%s#handleCreateResp() status_code_verification: != (201|400)", r.errCtx))
	}
}

func (r httpRepository) parseSuccess(body io.ReadCloser) (*data, error) {
	var ret payload
	dec := json.NewDecoder(body)
	if err := r.decode(dec, &ret); err != nil {
		return nil, errors.Wrapf(err, "%s#parseSuccess() decode", r.errCtx)
	}

	return ret.Data, nil
}

func (r httpRepository) parseClientError(body io.ReadCloser) (*data, error) {
	type clientError struct {
		Message string `json:"error_message"`
	}

	var cr clientError
	dec := json.NewDecoder(body)
	if err := r.decode(dec, &cr); err != nil {
		return nil, errors.Wrapf(err, "%s#parseClientError() decode", r.errCtx)
	}

	return nil, errors.New(cr.Message)
}

func (r httpRepository) fetch(id string) (*data, error) {
	resp, err := r.request(_get, fmt.Sprintf("http://%s:%s/v1/organisation/accounts/%s", r.addr, r.port, id), nil)
	if err != nil {
		return nil, errors.Wrapf(err, "%s#fetch() request", r.errCtx)
	}
	defer resp.Body.Close()

	return r.handleFetchResp(resp)
}

func (r httpRepository) delete(id string, version int64) error {
	resp, err := r.request(_delete, fmt.Sprintf("http://%s:%s/v1/organisation/accounts/%s?version=%v", r.addr, r.port, id, version), nil)
	if err != nil {
		return errors.Wrapf(err, "%s#delete() request", r.errCtx)
	}
	defer resp.Body.Close()

	return nil
}

func (r httpRepository) handleFetchResp(resp *http.Response) (*data, error) {
	const (
		success     = 200
		clientError = 400
		notFound    = 404
	)

	switch resp.StatusCode {
	case success:
		return r.parseSuccess(resp.Body)
	case clientError:
		return r.parseClientError(resp.Body)
	case notFound:
		return nil, nil
	default:
		return nil, errors.New(fmt.Sprintf("%s#handleFetchResp() status_code_verification: != (200|40[04])", r.errCtx))
	}
}

func (r httpRepository) health() error {
	resp, err := r.request(_get, fmt.Sprintf("http://%s:%s/v1/health", r.addr, r.port), nil)
	if err != nil {
		return errors.Wrapf(err, "%s#health() request", r.errCtx)
	}
	defer resp.Body.Close()

	return nil
}

func (a addrOption) apply(opts *httpOptions) {
	opts.addr = string(a)
}
func (p portOption) apply(opts *httpOptions) {
	opts.port = string(p)
}
