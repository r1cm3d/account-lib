package acc

type (
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
