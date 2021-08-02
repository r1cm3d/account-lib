package acc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type payload struct {
	Data *account `json:"data"`
}

// Account represents an account in the form3 org section.
// See https://api-docs.form3.tech/api.html#organisation-accounts for
// more information about fields.
type account struct {
	Attributes     *attributes `json:"attributes,omitempty"`
	ID             string      `json:"id,omitempty"`
	OrganisationID string      `json:"organisation_id,omitempty"`
	Type           string      `json:"type,omitempty"`
	Version        *int64      `json:"version,omitempty"`
}

// Attributes represents account attributes in the form3 org section.
// See https://api-docs.form3.tech/api.html#organisation-accounts for
// more information about fields.
type attributes struct {
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

// FIXME: error 500
func create(acc account) error {
	data, err := json.Marshal(payload{Data: &acc})
	s := string(data)
	fmt.Println(s)

	if err != nil {
		return err
	}

	resp, err := http.Post("http://0.0.0.0:8080/v1/organisation/accounts", "application/json", bytes.NewBuffer(data))

	if err != nil {
		return err
	}

	var ret account

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
