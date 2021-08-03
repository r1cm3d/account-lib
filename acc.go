package acc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type (
	payload struct {
		Data *account `json:"data"`
	}

	account struct {
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

func create(acc account) error {
	data, err := json.Marshal(payload{Data: &acc})

	if err != nil {
		return err
	}

	// TODO: change to use configuration
	resp, err := http.Post("http://0.0.0.0:8080/v1/organisation/accounts", "application/json", bytes.NewBuffer(data))

	if err != nil {
		return err
	}

	var ret account

	// TODO: Implement the unmarshalling logic here
	json.NewDecoder(resp.Body).Decode(&ret)

	fmt.Println(ret)
	return nil
}

// TODO: improve it
func health() error {
	resp, err := http.Get("http://0.0.0.0:8080/v1/health")

	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var data map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&data)

	fmt.Println(data)
	return nil
}
