package account

import (
	"errors"
	"github.com/google/uuid"
	"reflect"
	"testing"
)

func TestCreate(t *testing.T) {
	cases := []struct {
		name string
		in   Service
		want *Entity
	}{
		{"fully filled", serviceWithMockedRepository, filledEntityStub()},
	}
	cr := CreateRequest{
		OrganisationID:          _organisationIDStub,
		Number:                  _numberStub,
		AlternativeNames:        _alternativeNamesStub,
		BankID:                  _bankIDStub,
		BankIDCode:              _bankIDCodeStub,
		BaseCurrency:            _baseCurrencyStub,
		Bic:                     _bicStub,
		Country:                 _countryStub,
		Iban:                    _ibanStub,
		JointAccount:            _jointAccountStub,
		Name:                    _nameStub,
		SecondaryIdentification: _secondaryIdentificationStub,
		Switched:                _switchedStub,
		MatchingOptOut:          _matchingOptOutStub,
		Classification:          _classificationStub,
	}

	for _, tt := range cases {
		got, _ := tt.in.Create(cr)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Create(%v) got: %v, want: %v", tt.name, got, tt.want)
		}
	}
}

func TestCreate_Error(t *testing.T) {
	cases := []struct {
		name string
		in   Service
		want error
	}{
		{"repo create error", serviceWithCreateError, errors.New("service create_repo_create: organisationID: , country: : repo create error")},
		{"ofAcc error", serviceWithOutputMapperError, errors.New("service create_ofAcc: organisationID: , country: : ofAcc error")},
	}
	cr := CreateRequest{}

	for _, tt := range cases {
		_, got := tt.in.Create(cr)
		if got.Error() != tt.want.Error() {
			t.Errorf("Create_Error(%v) got: %v, want: %v", tt.name, got, tt.want)
		}
	}
}

const (
	_idStub                      = "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"
	_organisationIDStub          = "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c"
	_numberStub                  = "666"
	_bankIDStub                  = "400300"
	_bankIDCodeStub              = "GBDSC"
	_baseCurrencyStub            = "GBP"
	_bicStub                     = "NWBKGB22"
	_ibanStub                    = "GB33BUKB20201555555555"
	_secondaryIdentificationStub = "20530441"
)

var (
	_versionStub             = int64(0)
	_uuidStub, _             = uuid.Parse(_idStub)
	_organisationUUIDStub, _ = uuid.Parse(_organisationIDStub)
	_alternativeNamesStub    = []string{"Adanedhel"}
	_nameStub                = []string{"TURIN TURAMBAR"}
	_classificationStub      = "Personal"
	_matchingOptOutStub      = true
	_countryStub             = "GB"
	_jointAccountStub        = true
	_statusStub              = "confirmed"
	_switchedStub            = true
)

var (
	serviceWithCreateError = Service{
		errCtx:      "service",
		inputMapper: mockInputMapper{},
		creator:     mockCreatorErr{},
	}
	serviceWithOutputMapperError = Service{
		errCtx:       "service",
		inputMapper:  mockInputMapper{},
		creator:      mockCreatorOk{},
		outputMapper: mockOutputMapperErr{},
	}
	serviceWithMockedRepository = Service{
		errCtx:       "service",
		inputMapper:  mapper{},
		outputMapper: mapper{},
		creator:      mockCreatorOk{assertArg: true},
	}
)

func filledEntityStub() *Entity {
	return &Entity{
		id:                      _uuidStub,
		version:                 _versionStub,
		organisationID:          _organisationUUIDStub,
		classification:          Classification(_classificationStub),
		matchingOptOut:          _matchingOptOutStub,
		number:                  _numberStub,
		alternativeNames:        _alternativeNamesStub,
		bankID:                  _bankIDStub,
		bankIDCode:              _bankIDCodeStub,
		baseCurrency:            _baseCurrencyStub,
		bic:                     _bicStub,
		country:                 Country(_countryStub),
		iban:                    _ibanStub,
		jointAccount:            _jointAccountStub,
		name:                    _nameStub,
		secondaryIdentification: _secondaryIdentificationStub,
		status:                  Status(_statusStub),
		switched:                _switchedStub,
	}
}

type (
	mockCreatorErr struct{}
	mockCreatorOk  struct {
		assertArg bool
	}
	mockInputMapper     struct{}
	mockOutputMapperErr struct{}
)

func (m mockCreatorErr) create(_ data) (*data, error) {
	return nil, errors.New("repo create error")
}

func (m mockCreatorOk) create(d data) (*data, error) {
	expInput := &data{
		Attributes: &attributes{
			Number:                  _numberStub,
			AlternativeNames:        _alternativeNamesStub,
			Classification:          &_classificationStub,
			MatchingOptOut:          &_matchingOptOutStub,
			BankID:                  _bankIDStub,
			BankIDCode:              _bankIDCodeStub,
			BaseCurrency:            _baseCurrencyStub,
			Bic:                     _bicStub,
			Country:                 &_countryStub,
			Iban:                    _ibanStub,
			JointAccount:            &_jointAccountStub,
			Name:                    _nameStub,
			SecondaryIdentification: _secondaryIdentificationStub,
			Switched:                &_switchedStub,
		},
		OrganisationID: _organisationIDStub,
		Type:           "accounts",
		Version:        &_versionStub,
	}

	if m.assertArg &&
		(expInput.ID != d.ID ||
			expInput.Type != d.Type ||
			expInput.OrganisationID != d.OrganisationID ||
			!reflect.DeepEqual(expInput.Version, d.Version) ||
			!reflect.DeepEqual(expInput.Attributes, d.Attributes)) {
		return nil, errors.New("mockCreatorOk expInput did not match with data")
	}

	return &data{
		Attributes: &attributes{
			Classification:          &_classificationStub,
			MatchingOptOut:          &_matchingOptOutStub,
			Number:                  _numberStub,
			AlternativeNames:        _alternativeNamesStub,
			BankID:                  _bankIDStub,
			BankIDCode:              _bankIDCodeStub,
			BaseCurrency:            _baseCurrencyStub,
			Bic:                     _bicStub,
			Country:                 &_countryStub,
			Iban:                    _ibanStub,
			JointAccount:            &_jointAccountStub,
			Name:                    _nameStub,
			SecondaryIdentification: _secondaryIdentificationStub,
			Status:                  &_statusStub,
			Switched:                &_switchedStub,
		},
		ID:             _idStub,
		OrganisationID: _organisationIDStub,
		Type:           "accounts",
		Version:        &_versionStub,
	}, nil
}

func (m mockInputMapper) toAcc(_ CreateRequest) *data {
	return &data{}
}

func (m mockOutputMapperErr) ofAcc(data) (*Entity, error) {
	return nil, errors.New("ofAcc error")
}
