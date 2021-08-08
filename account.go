package account

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type (
	// Classification of account, only used for Confirmation of Payee (CoP).
	//
	// CoP: Can be either Personal or Business. Defaults to Personal if not provided.
	Classification string
	// Status of the account.
	//
	// FPS: Can be pending, confirmed or closed. (ALWAYS)
	//
	// SEPA & FPS Indirect (LHV): Can be either pending, confirmed or failed. (ALWAYS)
	//
	// All other services: Can be pending or confirmed. pending is a virtual state and is immediately superseded by confirmed. (ALWAYS)
	Status string
	// Currency refers to ISO 4217 code used to identify the base currency of the account, e.g. 'GBP', 'EUR'.
	//
	// See: https://www.iso.org/iso-4217-currency-codes.html
	Currency string
	// Country refers to ISO 3166-1 code used to identify the domicile of the account, e.g. 'GB', 'FR'.
	//
	// See: https://www.iso.org/iso-3166-country-codes.html
	Country string
	// CreateRequest groups attributes that are involved when creating an Account resource.
	//
	// See: https://api-docs.form3.tech/api.html#organisation-accounts
	CreateRequest struct {
		// ID is the unique ID of the resource in UUID 4 format. It identifies the resource within the system.
		//
		// Must be a new unique UUID 4 that hasn't been used in the Form3 system before. The call will fail with a 409
		// HTTP error code if a duplicate UUID is used. (REQUIRED)
		//
		// See: https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)
		ID string
		// OrganisationID of the organisation by which this resource has been created.
		//
		// Must be your organisation ID
		OrganisationID string
		// Classification is the classification of the account. (REQUIRED)
		Classification string
		// MatchingOptOut is a flag to indicate if the account has opted out of account matching, only used for
		// Confirmation of Payee. (OPTIONAL)
		//
		// CoP: Set to true if the account has opted out of account matching. Defaults to false.
		MatchingOptOut bool
		// Number is the unique account number. It will automatically be generated if not provided. If provided, the
		// account number is not validated. (OPTIONAL)
		Number string
		// AlternativeNames refers to the primary account names, only used for UK Confirmation of Payee. (OPTIONAL)
		//
		// CoP: Up to 3 alternative account names, one in each line of the array.
		AlternativeNames []string
		// BankID refers to local country bank identifier. Format depends on the country. Required for most
		// countries. (OPTIONAL)
		BankID string
		// BankIDCode identifies the type of bank ID being used. Required value depends on country attribute. (OPTIONAL)
		//
		// See: https://api-docs.form3.tech/api.html#accounts-create-data-table
		BankIDCode string
		// BaseCurrency is the Currency of the account. (CONDITIONAL)
		BaseCurrency string
		// Bic refers to the SWIFT BIC in either 8 or 11 character format e.g. 'NWBKGB22' (OPTIONAL)
		Bic string
		// Country refers to Country of the account. (OPTIONAL)
		Country string
		// Iban of the account. Will be calculated from other fields if not supplied. Ignored in SEPA Indirect,
		// provided by LHV after account generation is successful. (REQUIRED)
		Iban string
		// JointAccount is a flag to indicate if the account is a joint account, only used for Confirmation of Payee (CoP)
		//
		// CoP: Set to true is this is a joint account. Defaults to false if not provided. (OPTIONAL)
		JointAccount bool
		// Name of the account holder, up to four lines possible.
		//
		// CoP: Primary account name. For concatenated personal names, joint account names and organisation names,
		// use the first line. If first and last names of a personal name are separated, use the first line for first
		// names, the second line for last names. Titles are ignored and should not be entered. (REQUIRED)
		//
		// SEPA Indirect: Can be a person or organisation. Only the first line is used, minimum 5 characters. (REQUIRED)
		Name []string
		// SecondaryIdentification is the additional information to identify the account and account holder, only used
		// for Confirmation of Payee (CoP).
		//
		// CoP: Can be any type of additional identification, e.g. a building society roll number (OPTIONAL)
		SecondaryIdentification string
		// Switched is a flag to indicate if the account has been switched away from this organisation, only used for
		// Confirmation of Payee (CoP).
		//
		// CoP: Set to true if the account has been switched using the Current Account Switching Service (CASS),
		// false otherwise. (OPTIONAL)
		Switched bool
	}
	// Entity provides an abstraction to account. All information are provided by get methods
	Entity struct {
		id                      uuid.UUID
		version                 int64
		organisationID          uuid.UUID
		classification          Classification
		matchingOptOut          bool
		number                  string
		alternativeNames        []string
		bankID                  string
		bankIDCode              string
		baseCurrency            Currency
		bic                     string
		country                 Country
		iban                    string
		jointAccount            bool
		name                    []string
		secondaryIdentification string
		status                  Status
		switched                bool
	}
	// Service provides the main API to interact with account-api.
	//
	// It should not be instantiate directly. Use NewService(repo repository) *Service instead.
	Service struct {
		inputMapper
		outputMapper
		creator

		errCtx string
	}
)

type (
	mapper     struct{}
	repository interface {
		creator
	}
	creator interface {
		create(data) (*data, error)
	}
	inputMapper interface {
		toAcc(CreateRequest) *data
	}
	outputMapper interface {
		ofAcc(data) (*Entity, error)
	}
)

// NewService instantiates a Service. It is the only way to instantiate Service.
//
// It receives a repository as argument. The argument provides low level RPC to interact with account-api.
func NewService(repo repository) *Service {
	mapper := mapper{}
	return &Service{
		errCtx:       "service",
		creator:      repo,
		inputMapper:  mapper,
		outputMapper: mapper,
	}
}

// Create registers an existing bank account with account-api or create a new one. The Country attribute must be
// specified as a minimum. Depending on the country, other attributes such as BankID and Bic are mandatory.
//
// Returns error when CreateRequest -> data, repo.create(), data -> Entity fails.
func (s Service) Create(cr CreateRequest) (*Entity, error) {
	wrapErr := func(err error, msg string) error {
		return errors.Wrapf(err, "%s create_%s: organisationID: %s, country: %s", s.errCtx, msg, cr.OrganisationID, cr.Country)
	}
	data := s.toAcc(cr)

	ret, err := s.create(*data)
	if err != nil {
		return nil, wrapErr(err, "repo_create")
	}

	acc, err := s.ofAcc(*ret)
	if err != nil {
		return nil, wrapErr(err, "ofAcc")
	}

	return acc, nil
}

func (a Entity) ID() uuid.UUID {
	return a.id
}

func (a Entity) Version() int64 {
	return a.version
}

func (a Entity) OrganisationID() uuid.UUID {
	return a.organisationID
}

func (a Entity) Classification() Classification {
	return a.classification
}

func (a Entity) MatchingOptOut() bool {
	return a.matchingOptOut
}

func (a Entity) Number() string {
	return a.number
}

func (a Entity) AlternativeNames() []string {
	newAltNam := make([]string, len(a.alternativeNames))
	copy(newAltNam, a.alternativeNames)

	return newAltNam
}

func (a Entity) BankID() string {
	return a.bankID
}

func (a Entity) BankIDCode() string {
	return a.bankIDCode
}

func (a Entity) BaseCurrency() Currency {
	return a.baseCurrency
}

func (a Entity) Bic() string {
	return a.bic
}

func (a Entity) Country() Country {
	return a.country
}

func (a Entity) Iban() string {
	return a.iban
}

func (a Entity) JointAccount() bool {
	return a.jointAccount
}

func (a Entity) Name() []string {
	newName := make([]string, len(a.name))
	copy(newName, a.name)

	return newName
}

func (a Entity) SecondaryIdentification() string {
	return a.secondaryIdentification
}

func (a Entity) Status() Status {
	return a.status
}

func (a Entity) Switched() bool {
	return a.switched
}

func (r mapper) toAcc(cr CreateRequest) *data {
	defaultVersion := int64(0)
	return &data{
		Attributes: &attributes{
			Classification:          &cr.Classification,
			MatchingOptOut:          &cr.MatchingOptOut,
			Number:                  cr.Number,
			AlternativeNames:        cr.AlternativeNames,
			BankID:                  cr.BankID,
			BankIDCode:              cr.BankIDCode,
			BaseCurrency:            cr.BaseCurrency,
			Bic:                     cr.Bic,
			Country:                 &cr.Country,
			Iban:                    cr.Iban,
			JointAccount:            &cr.JointAccount,
			Name:                    cr.Name,
			SecondaryIdentification: cr.SecondaryIdentification,
			Switched:                &cr.Switched,
		},
		OrganisationID: cr.OrganisationID,
		Type:           "accounts",
		Version:        &defaultVersion,
		ID:             cr.ID,
	}
}

func (r mapper) ofAcc(d data) (*Entity, error) {
	id, err := uuid.Parse(d.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "ID parse: %s", d.ID)
	}

	organisationID, err := uuid.Parse(d.OrganisationID)
	if err != nil {
		return nil, errors.Wrapf(err, "organisationID parse: %s", d.OrganisationID)
	}

	att := d.Attributes
	if att == nil {
		return nil, errors.New("att.Attributes is nil")
	}

	var version int64
	if d.Version != nil {
		version = *d.Version
	}

	var classification Classification
	if att.Classification != nil {
		classification = Classification(*att.Classification)
	}

	var matchingOptOut bool
	if att.MatchingOptOut != nil {
		matchingOptOut = *att.MatchingOptOut
	}

	var country Country
	if att.Country != nil {
		country = Country(*att.Country)
	}

	var jointAccount bool
	if att.JointAccount != nil {
		jointAccount = *att.JointAccount
	}

	var status Status
	if att.Status != nil {
		status = Status(*att.Status)
	}

	var switched bool
	if att.Switched != nil {
		switched = *att.Switched
	}

	return &Entity{
		id:                      id,
		version:                 version,
		organisationID:          organisationID,
		classification:          classification,
		matchingOptOut:          matchingOptOut,
		number:                  att.Number,
		alternativeNames:        att.AlternativeNames,
		bankID:                  att.BankID,
		bankIDCode:              att.BankIDCode,
		baseCurrency:            Currency(att.BaseCurrency),
		bic:                     att.Bic,
		country:                 country,
		iban:                    att.Iban,
		jointAccount:            jointAccount,
		name:                    att.Name,
		secondaryIdentification: att.SecondaryIdentification,
		status:                  status,
		switched:                switched,
	}, nil
}
