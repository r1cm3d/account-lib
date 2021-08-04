package acc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
	"testing"
)

var (
	_fakeStubID = "ad27e265-9605-4b4b-a0e5-3003ea9cc4d2"
	_itAddress  = "0.0.0.0"
	_itPort     = "8080"
)

var (
	repoWithMarshalError = repo{
		marshal: func(v interface{}) ([]byte, error) { return nil, errors.New("error on marshal") },
	}
	repoWithPostError = repo{
		marshal: func(v interface{}) ([]byte, error) { return []byte("mock"), nil },
		post:    func(url, contentType string, body io.Reader) (resp *http.Response, err error) { return nil, errors.New("error on post") },
	}
	repoWithUnsuccessfullyStatusCode = repo{
		marshal: func(v interface{}) ([]byte, error) { return []byte("mock"), nil },
		post: func(url, contentType string, body io.Reader) (resp *http.Response, err error) {
			return &http.Response{
				StatusCode: 404,
			}, nil
		},
	}
	repoWithDecodeError = repo{
		marshal: func(v interface{}) ([]byte, error) { return []byte("mock"), nil },
		post: func(url, contentType string, body io.Reader) (resp *http.Response, err error) {
			return &http.Response{
				StatusCode: 201,
			}, nil
		},
		decode: func(d *json.Decoder, v interface{}) error { return errors.New("error on decode") },
	}
)

func TestCreateIntegration(t *testing.T) {
	skipShort(t)
	deleteStub(t)
	account := stubAccount()
	repo := newRepo(_itAddress, _itPort)

	got, err := repo.create(account)
	if err != nil {
		t.Fail()
	}

	fmt.Printf("Created account: %v", got.ID)
}

func TestCreate_Error(t *testing.T) {
	cases := []struct {
		name string
		in   repo
		want error
	}{
		{"marshal error", repoWithMarshalError, errors.New("http_repo create marshal: error on marshal")},
		{"post error", repoWithPostError, errors.New("http_repo create request: error on post")},
		{"unsuccessfully status code", repoWithUnsuccessfullyStatusCode, errors.New("http_repo create status code verification: not success != 201")},
		{"decode error", repoWithDecodeError, errors.New("http_repo create decode: error on decode")},
	}
	acc := stubAccount()

	for _, tt := range cases {
		_, got := tt.in.create(acc)
		if got.Error() != tt.want.Error() {
			t.Errorf("Create_Error(%v) got: %v, want: %v", tt.name, got, tt.want)
		}
	}
}

// TODO: improve it
func TestHealth(t *testing.T) {
	skipShort(t)
	repo := newRepo(_itAddress, _itPort)

	if err := repo.health(); err != nil {
		t.Fail()
	}
}

func skipShort(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
}

func stubAccount() account {
	country := "GB"

	return account{
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

	req, err := http.NewRequest("DELETE", fmt.Sprintf("http://%s:%s/v1/organisation/accounts/%s?version=0", _itAddress, _itPort, _fakeStubID), nil)
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

	resp, err := http.Post(fmt.Sprintf("http://%s:%s/v1/organisation/accounts", _itAddress, _itPort), "application/json", bytes.NewBuffer(data))

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != success {
		t.Fatal("Error on addStub")
	}
}
