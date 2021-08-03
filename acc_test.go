package acc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"
)

var (
	_fakeStubID = "ad27e265-9605-4b4b-a0e5-3003ea9cc4d2"
	_itAddress = "0.0.0.0"
	_itPort = "8080"
)

func TestCreateIntegration(t *testing.T) {
	skipShort(t)
	deleteStub(t)
	// TODO: test passing all arguments

	// WITH CoP
	// WITHOUT CoP
	account := stubAccount()

	if err := create(account); err != nil {
		t.Fail()
	}
}

// TODO: improve it
func TestHealth(t *testing.T) {
	skipShort(t)

	if err := health(); err != nil {
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
		success = 204
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

	fmt.Println(resp)
	fmt.Println(resp.Status)
	fmt.Println(resp.StatusCode)
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