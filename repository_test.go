package account

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type mockCloser struct {
	io.Reader
}

var (
	_fakeStubID      = "ad27e265-9605-4b4b-a0e5-3003ea9cc4d2"
	_fakeStubUUID, _ = uuid.Parse(_fakeStubID)
	_itAddress       = flag.String("itaddr", "0.0.0.0", "address of account-api service")
	_itPort          = "8080"
)

var (
	_mockedBytes                = []byte("mock")
	_repositoryWithMarshalError = httpRepository{
		errCtx:  "http_repository",
		marshal: func(v interface{}) ([]byte, error) { return nil, errors.New("error on marshal") },
	}
	_repositoryWithGetError = httpRepository{
		errCtx: "http_repository",
		get:    func(url string) (resp *http.Response, err error) { return nil, errors.New("error on get") },
	}
	_repositoryWithPostError = httpRepository{
		errCtx:  "http_repository",
		marshal: func(v interface{}) ([]byte, error) { return _mockedBytes, nil },
		post: func(url, contentType string, body io.Reader) (resp *http.Response, err error) {
			return nil, errors.New("error on post")
		},
	}
	_repositoryWithUnsuccessfullyStatusCode = httpRepository{
		errCtx:  "http_repository",
		marshal: func(v interface{}) ([]byte, error) { return _mockedBytes, nil },
		post: func(url, contentType string, body io.Reader) (resp *http.Response, err error) {
			return &http.Response{
				StatusCode: 404,
				Body:       mockCloser{bytes.NewBuffer(_mockedBytes)},
			}, nil
		},
	}
	_repositoryWithDecodeSuccessError = httpRepository{
		errCtx:  "http_repository",
		marshal: func(v interface{}) ([]byte, error) { return _mockedBytes, nil },
		post: func(url, contentType string, body io.Reader) (resp *http.Response, err error) {
			return &http.Response{
				StatusCode: 201,
				Body:       mockCloser{bytes.NewBuffer(_mockedBytes)},
			}, nil
		},
		decode: func(d *json.Decoder, v interface{}) error { return errors.New("error on decode success") },
	}
	_repositoryWithDecodeBadRequestError = httpRepository{
		errCtx:  "http_repository",
		marshal: func(v interface{}) ([]byte, error) { return _mockedBytes, nil },
		post: func(url, contentType string, body io.Reader) (resp *http.Response, err error) {
			return &http.Response{
				StatusCode: 400,
				Body:       mockCloser{bytes.NewBuffer(_mockedBytes)},
			}, nil
		},
		decode: func(d *json.Decoder, v interface{}) error { return errors.New("error on decode badRequest") },
	}
)

var (
	_accountStub = data{
		Attributes: &attributes{
			BankID:       "400300",
			BankIDCode:   "GBDSC",
			BaseCurrency: "GBP",
			Bic:          "NWBKGB22",
			Country:      &_countryStub,
			Name:         []string{"BRUCE", "WAYNE"},
		},
		ID:             _fakeStubID,
		OrganisationID: "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
		Type:           "accounts",
		Version: &_versionStub,
	}
	_accountFailedStub = data{
		Type: "accounts",
	}
)

func TestRepositoryCreateIntegration(t *testing.T) {
	skipShort(t)
	deleteStub(t)
	account := _accountStub
	repo := NewHTTPRepository(WithAddr(*_itAddress), WithPort(_itPort))

	got, err := repo.create(account)
	if err != nil {
		t.Fatal()
	}

	fmt.Printf("Created data: %v", got.ID)
}

func TestRepositoryCreate_ErrorIntegration(t *testing.T) {
	skipShort(t)
	deleteStub(t)
	account := _accountFailedStub
	repo := NewHTTPRepository(WithAddr(*_itAddress), WithPort(_itPort))

	_, err := repo.create(account)
	if err == nil {
		t.Fatal()
	}

	fmt.Printf("Message error: %v", err.Error())
}

func TestHealthIntegration(t *testing.T) {
	skipShort(t)
	repo := NewHTTPRepository(WithAddr(*_itAddress), WithPort(_itPort))

	if err := repo.health(); err != nil {
		t.Fail()
	}
}

func TestFetchIntegration(t *testing.T) {
	errT := func(propName string, got, want interface{}) {
		t.Errorf("TestFetchIntegration_%s got: %v, want: %v", propName, got, want)
	}
	assertInterface := func(propName string, got, want interface{}) {
		if !reflect.DeepEqual(got, want) {
			errT(propName, got, want)
		}
	}
	assertData := func(propName string, got, want interface{}) {
		if !reflect.DeepEqual(got, want) {
			errT(propName, got, want)
		}
	}
	skipShort(t)
	addStub(t)
	repo := NewHTTPRepository(WithAddr(*_itAddress), WithPort(_itPort))

	if got, _ := repo.fetch(_fakeStubID); got != nil {
		assertInterface("Attribute", got.Attributes, _accountStub.Attributes)
		assertInterface("Version", got.Version, _accountStub.Version)
		assertData("ID", got.ID, _accountStub.ID)
		assertData("OrganisationID", got.OrganisationID, _accountStub.OrganisationID)
		assertData("OrganisationID", got.Type, _accountStub.Type)

		return
	}

	t.Fail()
}

func TestHealth_Error(t *testing.T) {
	cases := []struct {
		name string
		in   httpRepository
		want error
	}{
		{"get error", _repositoryWithGetError, errors.New("http_repository#health() get: error on get")},
	}

	for _, tt := range cases {
		got := tt.in.health()
		if got.Error() != tt.want.Error() {
			t.Errorf("RepositoryHealth_Error(%v) got: %v, want: %v", tt.name, got, tt.want)
		}
	}
}

func TestRepositoryCreate_Error(t *testing.T) {
	cases := []struct {
		name string
		in   httpRepository
		want error
	}{
		{"marshal", _repositoryWithMarshalError, errors.New("http_repository#create() marshal: error on marshal")},
		{"post", _repositoryWithPostError, errors.New("http_repository#create() request: error on post")},
		{"unsuccessfully status code", _repositoryWithUnsuccessfullyStatusCode, errors.New("http_repository#handleCreateResp() status_code_verification: != (201|400)")},
		{"decode success error", _repositoryWithDecodeSuccessError, errors.New("http_repository#parseSuccess() decode: error on decode success")},
		{"decode badRequest error", _repositoryWithDecodeBadRequestError, errors.New("http_repository#parseClientError() decode: error on decode badRequest")},
	}
	acc := _accountStub

	for _, tt := range cases {
		_, got := tt.in.create(acc)
		if got.Error() != tt.want.Error() {
			t.Errorf("RepositoryCreate_Error(%v) got: %v, want: %v", tt.name, got, tt.want)
		}
	}
}

func skipShort(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
}

func deleteStub(t *testing.T) {
	const (
		success  = 204
		notFound = 404
	)
	client := &http.Client{}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("http://%s:%s/v1/organisation/accounts/%s?version=0", *_itAddress, _itPort, _fakeStubID), nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if (resp.StatusCode != success) && (resp.StatusCode != notFound) {
		t.Fatalf("Error on delete %s", _fakeStubID)
	}
}

func addStub(t *testing.T) {
	deleteStub(t)
	const success = 201

	stub := _accountStub
	data, err := json.Marshal(payload{Data: &stub})

	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.Post(fmt.Sprintf("http://%s:%s/v1/organisation/accounts", *_itAddress, _itPort), "application/json", bytes.NewBuffer(data))

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != success {
		t.Fatal("Error on addStub")
	}
}

func (mockCloser) Close() error { return nil }
