package account

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func TestRepositoryCreateIntegration(t *testing.T) {
	skipShort(t)
	deleteStub(t)
	account := stubAccount()
	repo := NewHTTPRepository(WithAddr(*_itAddress), WithPort(_itPort))

	got, err := repo.create(account)
	if err != nil {
		t.Fatal()
	}

	fmt.Printf("Created data: %v", got.ID)
}

func TestRepositoryCreate_Error(t *testing.T) {
	cases := []struct {
		name string
		in   httpRepository
		want error
	}{
		{"marshal error", repositoryWithMarshalError, errors.New("http_repository create marshal: error on marshal")},
		{"post error", repositoryWithPostError, errors.New("http_repository create request: error on post")},
		{"unsuccessfully status code", repositoryWithUnsuccessfullyStatusCode, errors.New("http_repository create status code verification: not success != 201")},
		{"decode error", repositoryWithDecodeError, errors.New("http_repository create decode: error on decode")},
	}
	acc := stubAccount()

	for _, tt := range cases {
		_, got := tt.in.create(acc)
		if got.Error() != tt.want.Error() {
			t.Errorf("RepositoryCreate_Error(%v) got: %v, want: %v", tt.name, got, tt.want)
		}
	}
}

// TODO: improve it
func TestHealth(t *testing.T) {
	skipShort(t)
	repo := NewHTTPRepository(WithAddr(*_itAddress), WithPort(_itPort))

	if err := repo.health(); err != nil {
		t.Fail()
	}
}

func skipShort(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
}

func stubAccount() data {
	country := "GB"

	return data{
		Attributes: &attributes{
			BankID:       "400300",
			BankIDCode:   "GBDSC",
			BaseCurrency: "GBP",
			Bic:          "NWBKGB22",
			Country:      &country,
			Name:         []string{"BRUCE", "WAYNE"},
		},
		ID:             _fakeStubID,
		OrganisationID: "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
		Type:           "accounts",
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

	stub := stubAccount()
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

var (
	_fakeStubID      = "ad27e265-9605-4b4b-a0e5-3003ea9cc4d2"
	_fakeStubUUID, _ = uuid.Parse(_fakeStubID)
	_itAddress       = flag.String("itaddr", "0.0.0.0", "address of account-api service")
	_itPort          = "8080"
)

var (
	mockedBytes                = []byte("mock")
	repositoryWithMarshalError = httpRepository{
		errCtx:  "http_repository",
		marshal: func(v interface{}) ([]byte, error) { return nil, errors.New("error on marshal") },
	}
	repositoryWithPostError = httpRepository{
		errCtx:  "http_repository",
		marshal: func(v interface{}) ([]byte, error) { return mockedBytes, nil },
		post: func(url, contentType string, body io.Reader) (resp *http.Response, err error) {
			return nil, errors.New("error on post")
		},
	}
	repositoryWithUnsuccessfullyStatusCode = httpRepository{
		errCtx:  "http_repository",
		marshal: func(v interface{}) ([]byte, error) { return mockedBytes, nil },
		post: func(url, contentType string, body io.Reader) (resp *http.Response, err error) {
			return &http.Response{
				StatusCode: 404,
				Body:       mockCloser{bytes.NewBuffer(mockedBytes)},
			}, nil
		},
	}
	repositoryWithDecodeError = httpRepository{
		errCtx:  "http_repository",
		marshal: func(v interface{}) ([]byte, error) { return mockedBytes, nil },
		post: func(url, contentType string, body io.Reader) (resp *http.Response, err error) {
			return &http.Response{
				StatusCode: 201,
				Body:       mockCloser{bytes.NewBuffer(mockedBytes)},
			}, nil
		},
		decode: func(d *json.Decoder, v interface{}) error { return errors.New("error on decode") },
	}
)

type mockCloser struct {
	io.Reader
}

func (mockCloser) Close() error { return nil }
