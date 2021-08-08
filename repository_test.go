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

type (
	mockCloser struct {
		io.Reader
	}
	mockRequestError               struct{}
	mockRequestBadRequest          struct{}
	mockRequestOk                  struct{}
	mockRequestCreate              struct{}
	mockRequestNotFound            struct{}
	mockRequestInternalServerError struct{}
)

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
	_repositoryWithRequestError = httpRepository{
		errCtx: "http_repository",
		client: mockRequestError{},
	}
	_repositoryWithPostError = httpRepository{
		errCtx:  "http_repository",
		marshal: func(v interface{}) ([]byte, error) { return _mockedBytes, nil },
		client:  mockRequestError{},
	}
	_repositoryWithUnsuccessfullyStatusCodeCreate = httpRepository{
		errCtx:  "http_repository",
		marshal: func(v interface{}) ([]byte, error) { return _mockedBytes, nil },
		client:  mockRequestNotFound{},
	}
	_repositoryWithUnsuccessfullyStatusCode = httpRepository{
		errCtx:  "http_repository",
		marshal: func(v interface{}) ([]byte, error) { return _mockedBytes, nil },
		client:  mockRequestInternalServerError{},
	}
	_repositoryWithDecodeSuccessErrorCreate = httpRepository{
		errCtx:  "http_repository",
		marshal: func(v interface{}) ([]byte, error) { return _mockedBytes, nil },
		client:  mockRequestCreate{},
		decode:  func(d *json.Decoder, v interface{}) error { return errors.New("error on decode success") },
	}
	_repositoryWithDecodeSuccessErrorFetch = httpRepository{
		errCtx:  "http_repository",
		marshal: func(v interface{}) ([]byte, error) { return _mockedBytes, nil },
		client:  mockRequestOk{},
		decode:  func(d *json.Decoder, v interface{}) error { return errors.New("error on decode success") },
	}
	_repositoryWithDecodeBadRequestError = httpRepository{
		errCtx:  "http_repository",
		marshal: func(v interface{}) ([]byte, error) { return _mockedBytes, nil },
		client:  mockRequestBadRequest{},
		decode:  func(d *json.Decoder, v interface{}) error { return errors.New("error on decode badRequest") },
	}
	_repositoryWithDecodeNotFoundError = httpRepository{
		errCtx:  "http_repository",
		marshal: func(v interface{}) ([]byte, error) { return _mockedBytes, nil },
		client:  mockRequestBadRequest{},
		decode:  func(d *json.Decoder, v interface{}) error { return errors.New("error on decode notFound") },
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
		Version:        &_versionStub,
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
		assertData("Type", got.Type, _accountStub.Type)

		return
	}

	t.Fail()
}

func TestDeleteIntegration(t *testing.T) {
	skipShort(t)
	addStub(t)
	repo := NewHTTPRepository(WithAddr(*_itAddress), WithPort(_itPort))

	if err := repo.delete(_fakeStubID, int64(0)); err != nil {
		t.Fail()
	}
}

func TestFetch_NotFoundIntegration(t *testing.T) {
	skipShort(t)
	addStub(t)
	repo := NewHTTPRepository(WithAddr(*_itAddress), WithPort(_itPort))
	nfID := "eed9954a-a58b-4f59-b44f-8d0592748d53"

	got, err := repo.fetch(nfID)

	if got != nil || err != nil {
		t.Fail()
	}
}

func TestRepositoryFetch_ErrorIntegration(t *testing.T) {
	skipShort(t)
	repo := NewHTTPRepository(WithAddr(*_itAddress), WithPort(_itPort))
	invalidID := "666"

	_, err := repo.fetch(invalidID)
	if err == nil {
		t.Errorf("RepositoryFetch_ErrorIntegration in: %v, want: NOT ERROR", invalidID)
	}

	fmt.Printf("Message error: %v", err.Error())
}

func TestHealth_Error(t *testing.T) {
	error := "http_repository#health() request: error on request"

	got := _repositoryWithRequestError.health()

	if got.Error() != error {
		t.Errorf("RepositoryHealth_Error got: %v, want: %v", got, error)
	}
}

func TestRepositoryCreate_Error(t *testing.T) {
	cases := []struct {
		name string
		in   httpRepository
		want error
	}{
		{"marshal", _repositoryWithMarshalError, errors.New("http_repository#create() marshal: error on marshal")},
		{"post", _repositoryWithPostError, errors.New("http_repository#create() request: error on request")},
		{"unsuccessfully status code", _repositoryWithUnsuccessfullyStatusCode, errors.New("http_repository#handleCreateResp() status_code_verification: != (201|400)")},
		{"decode success error", _repositoryWithDecodeSuccessErrorCreate, errors.New("http_repository#parseSuccess() decode: error on decode success")},
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

func TestRepositoryFetch_Error(t *testing.T) {
	cases := []struct {
		name string
		in   httpRepository
		want error
	}{
		{"request", _repositoryWithRequestError, errors.New("http_repository#fetch() request: error on request")},
		{"unsuccessfully status code", _repositoryWithUnsuccessfullyStatusCode, errors.New("http_repository#handleFetchResp() status_code_verification: != (200|40[04])")},
		{"decode success error", _repositoryWithDecodeSuccessErrorFetch, errors.New("http_repository#parseSuccess() decode: error on decode success")},
		{"decode badRequest error", _repositoryWithDecodeBadRequestError, errors.New("http_repository#parseClientError() decode: error on decode badRequest")},
	}
	for _, tt := range cases {
		_, got := tt.in.fetch(_fakeStubID)
		if got.Error() != tt.want.Error() {
			t.Errorf("RepositoryFetch_Error(%v) got: %v, want: %v", tt.name, got, tt.want)
		}
	}
}

func TestRepositoryDelete_Error(t *testing.T) {
	want := "http_repository#delete() request: error on request"

	got := _repositoryWithRequestError.delete(_fakeStubID, int64(0))
	if got.Error() != want {
		t.Errorf("RepositoryDelete_Error got: %v, want: %v", got, want)
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

func (r mockRequestError) request(_ method, _ string, _ io.Reader) (*http.Response, error) {
	return nil, errors.New("error on request")
}

func (r mockRequestBadRequest) request(_ method, _ string, _ io.Reader) (*http.Response, error) {
	return &http.Response{
		StatusCode: 400,
		Body:       mockCloser{bytes.NewBuffer(_mockedBytes)},
	}, nil
}

func (r mockRequestOk) request(_ method, _ string, _ io.Reader) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Body:       mockCloser{bytes.NewBuffer(_mockedBytes)},
	}, nil
}

func (r mockRequestCreate) request(_ method, _ string, _ io.Reader) (*http.Response, error) {
	return &http.Response{
		StatusCode: 201,
		Body:       mockCloser{bytes.NewBuffer(_mockedBytes)},
	}, nil
}

func (r mockRequestNotFound) request(_ method, _ string, _ io.Reader) (*http.Response, error) {
	return &http.Response{
		StatusCode: 404,
		Body:       mockCloser{bytes.NewBuffer(_mockedBytes)},
	}, nil
}

func (r mockRequestInternalServerError) request(_ method, _ string, _ io.Reader) (*http.Response, error) {
	return &http.Response{
		StatusCode: 500,
		Body:       mockCloser{bytes.NewBuffer(_mockedBytes)},
	}, nil
}
