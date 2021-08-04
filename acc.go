package acc

import "github.com/pkg/errors"

type (
	CreateRequest struct {
	}
	Account struct {
	}
	Service struct {
		toAcc
		ofAcc
		creator

		errCtx string
	}
	repository interface {
		creator
	}
	creator interface {
		create(data) (*data, error)
	}
	toAcc func(CreateRequest) (*data, error)
	ofAcc func(data) (*Account, error)
)

// TODO: Create an integration test for all of it
func NewService(repo repository) *Service {
	return &Service{
		toAcc:   nil,
		ofAcc:   nil,
		errCtx:  "service",
		creator: repo,
	}
}

func (s Service) Create(cr CreateRequest) (*Account, error) {
	wrapErr := func(err error, msg string) error {
		return errors.Wrapf(err, "%s create_%s: %v", s.errCtx, msg, cr)
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
