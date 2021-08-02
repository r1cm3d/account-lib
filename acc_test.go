package acc

import "testing"

func TestCreateIntegration(t *testing.T) {
	skipShort(t)
	// TODO: Prepare table here
	country := "GB"
	// TODO: test passing all arguments

	// WITH CoP
	// WITHOUT CoP
	account := account{
		Attributes: &attributes{
			BankID:                 "400300",
			BankIDCode:             "GBDSC",
			BaseCurrency:           "GBP",
			Bic:                    "NWBKGB22",
			Country:                &country,
			Name:                   []string{"BRUCE", "WAYNE"},
		},
		ID:             "ad27e265-9605-4b4b-a0e5-3003ea9cc4d2",
		OrganisationID: "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
		Type:           "accounts",
	}

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
