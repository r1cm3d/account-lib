package account

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

type (
	mockErr             struct{}
	mockInputMapper     struct{}
	mockOutputMapperErr struct{}
	mockOk              struct {
		assertArg bool
		expData   data
	}
)

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
	_serviceWithRequestError = Service{
		errCtx:      "service",
		inputMapper: mockInputMapper{},
		creator:     mockErr{},
		retriever:   mockErr{},
		eraser:      mockErr{},
	}
	_serviceWithOutputMapperError = Service{
		errCtx:       "service",
		inputMapper:  mockInputMapper{},
		creator:      mockOk{expData: _fullyFilledData},
		retriever:    mockOk{expData: _fullyFilledData},
		outputMapper: mockOutputMapperErr{},
	}
	_serviceWithMockedRepositoryFullyFilled = Service{
		errCtx:       "service",
		inputMapper:  mapper{},
		outputMapper: mapper{},
		creator:      mockOk{assertArg: true, expData: _fullyFilledData},
		retriever:    mockOk{assertArg: true, expData: _fullyFilledData},
	}
	_serviceWithMockedRepositoryBasicFilled = Service{
		errCtx:       "service",
		inputMapper:  mapper{},
		outputMapper: mapper{},
		creator:      mockOk{assertArg: true, expData: _basicFilledData},
	}
)

var (
	_dataWithIDErr = data{
		ID: "3rr0r",
	}
	_dataWithOrganizationIDErr = data{
		ID:             _idStub,
		OrganisationID: "0RG4N1Z4T10N_3rr0r",
	}
	_dataWithAttErr = data{
		Attributes:     nil,
		ID:             _idStub,
		OrganisationID: _organisationIDStub,
	}
)

var (
	_fullyFilledCreateRequest = CreateRequest{
		ID:                      _idStub,
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
	_fullyFilledEntity = &Entity{
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
	_fullyFilledData = data{
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
		OrganisationID: _organisationIDStub,
		Type:           "accounts",
		Version:        &_versionStub,
		ID:             _idStub,
	}
	_basicFilledCreateRequest = CreateRequest{
		ID:             _fakeStubID,
		OrganisationID: _organisationIDStub,
		Classification: _classificationStub,
		Number:         _numberStub,
		BankID:         _bankIDStub,
		BankIDCode:     _bankIDCodeStub,
		BaseCurrency:   _baseCurrencyStub,
		Bic:            _bicStub,
		Country:        _countryStub,
		Iban:           _ibanStub,
		Name:           _nameStub,
	}
	_basicFilledEntity = &Entity{
		id:             _fakeStubUUID,
		version:        _versionStub,
		organisationID: _organisationUUIDStub,
		classification: Classification(_classificationStub),
		number:         _numberStub,
		bankID:         _bankIDStub,
		bankIDCode:     _bankIDCodeStub,
		baseCurrency:   _baseCurrencyStub,
		bic:            _bicStub,
		country:        Country(_countryStub),
		iban:           _ibanStub,
		name:           _nameStub,
	}
	_basicMatchingOptOutStub = false
	_basicJointAccountStub   = false
	_basicSwitchedStub       = false
	_basicFilledData         = data{
		Attributes: &attributes{
			Number:         _numberStub,
			BankID:         _bankIDStub,
			BankIDCode:     _bankIDCodeStub,
			BaseCurrency:   _baseCurrencyStub,
			Classification: &_classificationStub,
			Bic:            _bicStub,
			Country:        &_countryStub,
			Iban:           _ibanStub,
			Name:           _nameStub,
			MatchingOptOut: &_basicMatchingOptOutStub,
			JointAccount:   &_basicJointAccountStub,
			Switched:       &_basicSwitchedStub,
		},
		ID:             _fakeStubID,
		OrganisationID: _organisationIDStub,
		Type:           "accounts",
		Version:        &_versionStub,
	}
)

func TestAccountCreateIntegration(t *testing.T) {
	skipShort(t)
	deleteStub(t)

	svc := NewService(NewHTTPRepository(WithAddr(*_itAddress)))

	got, err := svc.Create(_basicFilledCreateRequest)
	if err != nil {
		t.Fatal()
	}

	fmt.Printf("Created data: %v", got.ID())
}

func TestAccountFetchIntegration(t *testing.T) {
	skipShort(t)
	addStub(t)

	svc := NewService(NewHTTPRepository(WithAddr(*_itAddress)))

	got, err := svc.Fetch(_idStub)
	if err != nil {
		t.Fatal()
	}

	fmt.Printf("Fetched data: %v", got.ID())
}

func TestAccountDeleteIntegration(t *testing.T) {
	skipShort(t)
	addStub(t)

	svc := NewService(NewHTTPRepository(WithAddr(*_itAddress)))

	err := svc.Delete(BuildDeleteRequest(_fakeStubID))
	if err != nil {
		t.Fatal()
	}

	fmt.Print("Deleted data")
}

func TestAccountCreate(t *testing.T) {
	type in struct {
		cr  CreateRequest
		svc Service
	}
	cases := []struct {
		name string
		in
		want *Entity
	}{
		{"fully filled", in{_fullyFilledCreateRequest, _serviceWithMockedRepositoryFullyFilled}, _fullyFilledEntity},
		{"basic filled", in{_basicFilledCreateRequest, _serviceWithMockedRepositoryBasicFilled}, _basicFilledEntity},
	}

	for _, tt := range cases {
		got, _ := tt.in.svc.Create(tt.in.cr)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Create(%v) got: %v, want: %v", tt.name, got, tt.want)
		}
	}
}

func TestAccountFetch(t *testing.T) {
	got, _ := _serviceWithMockedRepositoryFullyFilled.Fetch(_idStub)
	if !reflect.DeepEqual(got, _fullyFilledEntity) {
		t.Errorf("Fetch got: %v, want: %v", got, _fullyFilledEntity)
	}
}

func TestCreate_Error(t *testing.T) {
	cases := []struct {
		name string
		in   Service
		want error
	}{
		{"repo", _serviceWithRequestError, errors.New("service create_repo_create: organisationID: , country: : repo create error")},
		{"ofAcc", _serviceWithOutputMapperError, errors.New("service create_ofAcc: organisationID: , country: : ofAcc error")},
	}
	cr := CreateRequest{}

	for _, tt := range cases {
		_, got := tt.in.Create(cr)
		if got.Error() != tt.want.Error() {
			t.Errorf("Create_Error(%v) got: %v, want: %v", tt.name, got, tt.want)
		}
	}
}

func TestFetch_Error(t *testing.T) {
	cases := []struct {
		name string
		in   Service
		want error
	}{
		{"repo", _serviceWithRequestError, errors.New("service fetch_repo_fetch: id: id: repo fetch error")},
		{"ofAcc", _serviceWithOutputMapperError, errors.New("service fetch_ofAcc: id: id: ofAcc error")},
	}

	for _, tt := range cases {
		_, got := tt.in.Fetch("id")
		if got.Error() != tt.want.Error() {
			t.Errorf("Fetch_Error(%v) got: %v, want: %v", tt.name, got, tt.want)
		}
	}
}

func TestDelete_Error(t *testing.T) {
	msg := "service delete: id: ad27e265-9605-4b4b-a0e5-3003ea9cc4d2: repo delete error"
	if got := _serviceWithRequestError.Delete(BuildDeleteRequest(_fakeStubID)); got.Error() != msg {
		t.Errorf("Delete_Error got: %v, want: %v", got, msg)
	}
}

func TestOfAcc_Error(t *testing.T) {
	cases := []struct {
		name string
		in   data
		want error
	}{
		{"ID parser", _dataWithIDErr, errors.New("ID parse: 3rr0r: invalid UUID length: 5")},
		{"organisationID parser", _dataWithOrganizationIDErr, errors.New("organisationID parse: 0RG4N1Z4T10N_3rr0r: invalid UUID length: 18")},
		{"att nil", _dataWithAttErr, errors.New("att.Attributes is nil")},
	}
	m := mapper{}

	for _, tt := range cases {
		_, got := m.ofAcc(tt.in)
		if got.Error() != tt.want.Error() {
			t.Errorf("OfAcc_Error(%v) got: %v, want: %v", tt.name, got, tt.want)
		}
	}
}

func TestEntity(t *testing.T) {
	assert := func(propName string, got, want interface{}) {
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Entity_%s got: %v, want: %v", propName, got, want)
		}
	}
	entity := Entity{
		id:                      _uuidStub,
		version:                 _versionStub,
		organisationID:          _organisationUUIDStub,
		classification:          Classification(_classificationStub),
		matchingOptOut:          _matchingOptOutStub,
		number:                  _numberStub,
		alternativeNames:        _alternativeNamesStub,
		bankID:                  _bankIDStub,
		bankIDCode:              _bankIDCodeStub,
		baseCurrency:            Currency(_baseCurrencyStub),
		bic:                     _bicStub,
		country:                 Country(_countryStub),
		iban:                    _ibanStub,
		jointAccount:            _jointAccountStub,
		name:                    _nameStub,
		secondaryIdentification: _secondaryIdentificationStub,
		status:                  Status(_statusStub),
		switched:                _switchedStub,
	}

	assert("ID", entity.ID(), _idStub)
	assert("UUID", entity.UUID(), _uuidStub)
	assert("Version", entity.Version(), _versionStub)
	assert("OrganisationID", entity.OrganisationID(), _organisationUUIDStub)
	assert("Classification", entity.Classification(), Classification(_classificationStub))
	assert("MatchingOptOut", entity.MatchingOptOut(), _matchingOptOutStub)
	assert("Number", entity.Number(), _numberStub)
	assert("AlternativeNames", entity.AlternativeNames(), _alternativeNamesStub)
	assert("BankID", entity.BankID(), _bankIDStub)
	assert("BankIDCode", entity.BankIDCode(), _bankIDCodeStub)
	assert("BaseCurrency", entity.BaseCurrency(), Currency(_baseCurrencyStub))
	assert("Bic", entity.Bic(), _bicStub)
	assert("Country", entity.Country(), Country(_countryStub))
	assert("Iban", entity.Iban(), _ibanStub)
	assert("JointAccount", entity.JointAccount(), _jointAccountStub)
	assert("Name", entity.Name(), _nameStub)
	assert("SecondaryIdentification", entity.SecondaryIdentification(), _secondaryIdentificationStub)
	assert("Status", entity.Status(), Status(_statusStub))
	assert("Switched", entity.Switched(), _switchedStub)
}

func (m mockErr) create(_ data) (*data, error) {
	return nil, errors.New("repo create error")
}

func (m mockErr) fetch(_ string) (*data, error) {
	return nil, errors.New("repo fetch error")
}

func (m mockErr) delete(_ string, _ int64) error {
	return errors.New("repo delete error")
}

func (m mockOk) create(d data) (*data, error) {
	expInput := m.buildData()

	if m.assertArg &&
		(expInput.ID != d.ID ||
			expInput.Type != d.Type ||
			expInput.OrganisationID != d.OrganisationID ||
			!reflect.DeepEqual(expInput.Version, d.Version) ||
			!reflect.DeepEqual(expInput.Attributes, d.Attributes)) {
		return nil, errors.New("mockOk expInput did not match with data")
	}

	return &m.expData, nil
}

func (m mockOk) fetch(id string) (*data, error) {
	expInput := m.buildData()

	if m.assertArg && (expInput.ID != id) {
		return nil, errors.New("mockOk expInput did not match with data")
	}

	return &m.expData, nil
}

func (m mockOk) buildData() *data {
	return &data{
		Attributes: &attributes{
			Classification:          m.expData.Attributes.Classification,
			MatchingOptOut:          m.expData.Attributes.MatchingOptOut,
			Number:                  m.expData.Attributes.Number,
			AlternativeNames:        m.expData.Attributes.AlternativeNames,
			BankID:                  m.expData.Attributes.BankID,
			BankIDCode:              m.expData.Attributes.BankIDCode,
			BaseCurrency:            m.expData.Attributes.BaseCurrency,
			Bic:                     m.expData.Attributes.Bic,
			Country:                 m.expData.Attributes.Country,
			Iban:                    m.expData.Attributes.Iban,
			JointAccount:            m.expData.Attributes.JointAccount,
			Name:                    m.expData.Attributes.Name,
			SecondaryIdentification: m.expData.Attributes.SecondaryIdentification,
			Switched:                m.expData.Attributes.Switched,
		},
		OrganisationID: m.expData.OrganisationID,
		Type:           "accounts",
		Version:        m.expData.Version,
		ID:             m.expData.ID,
	}
}

func (m mockInputMapper) toAcc(_ CreateRequest) *data {
	return &data{}
}

func (m mockOutputMapperErr) ofAcc(data) (*Entity, error) {
	return nil, errors.New("ofAcc error")
}
