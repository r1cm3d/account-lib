package account

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type (
	Classification string
	Status         string
	Currency       string
	Country        string
	CreateRequest  struct {
		OrganisationID          string
		Classification          string
		MatchingOptOut          bool
		Number                  string
		AlternativeNames        []string
		BankID                  string
		BankIDCode              string
		BaseCurrency            string
		Bic                     string
		Country                 string
		Iban                    string
		JointAccount            bool
		Name                    []string
		SecondaryIdentification string
		Switched                bool
	}
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

	Service struct {
		inputMapper
		outputMapper
		creator

		errCtx string
	}
	mapper     struct{}
	repository interface {
		creator
	}
	creator interface {
		create(data) (*data, error)
	}
	inputMapper interface {
		toAcc(CreateRequest) (*data, error)
	}
	outputMapper interface {
		ofAcc(data) (*Entity, error)
	}
)

// TODO: Create an integration test for all of it. Maybe use a docker container
func NewService(repo repository) *Service {
	mapper := mapper{}
	return &Service{
		errCtx:       "service",
		creator:      repo,
		inputMapper:  mapper,
		outputMapper: mapper,
	}
}

func (s Service) Create(cr CreateRequest) (*Entity, error) {
	wrapErr := func(err error, msg string) error {
		return errors.Wrapf(err, "%s create_%s: organisationID: %s, country: %s", s.errCtx, msg, cr.OrganisationID, cr.Country)
	}
	data, err := s.toAcc(cr)
	if err != nil {
		return nil, wrapErr(err, "toAcc")
	}

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

func (r mapper) toAcc(cr CreateRequest) (*data, error) {
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
	}, nil
}

func (r mapper) ofAcc(d data) (*Entity, error) {
	id, err := uuid.Parse(d.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "id parse: %s", d.ID)
	}

	organisationID, err := uuid.Parse(d.OrganisationID)
	if err != nil {
		return nil, errors.Wrapf(err, "organisationID parse: %s", d.OrganisationID)
	}

	att := d.Attributes
	if att == nil {
		return nil, errors.Wrap(err, "Attributes is nil")
	}

	return &Entity{
		id:                      id,
		version:                 *d.Version,
		organisationID:          organisationID,
		classification:          Classification(*att.Classification),
		matchingOptOut:          *att.MatchingOptOut,
		number:                  att.Number,
		alternativeNames:        att.AlternativeNames,
		bankID:                  att.BankID,
		bankIDCode:              att.BankIDCode,
		baseCurrency:            Currency(att.BaseCurrency),
		bic:                     att.Bic,
		country:                 Country(*att.Country),
		iban:                    att.Iban,
		jointAccount:            *att.JointAccount,
		name:                    att.Name,
		secondaryIdentification: att.SecondaryIdentification,
		status:                  Status(*att.Status),
		switched:                *att.Switched,
	}, nil
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