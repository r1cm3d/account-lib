package acc

import "github.com/pkg/errors"

type (
	CreateRequest struct {
	}
	Account struct {
	}
	Service struct {
		toAcc
		creator

		errCtx string
	}
	repository interface {
		creator
	}
	creator interface {
		create(acc data) (*data, error)
	}
	toAcc func(CreateRequest) (*data, error)
)

// TODO: Create an integration test for all of it
func NewService(repo repository) *Service {
	return &Service{
		toAcc:   nil,
		errCtx:  "service",
		creator: repo,
	}
}

func (s Service) Create(cr CreateRequest) (*Account, error) {
	data, err := s.toAcc(cr)
	if err != nil {
		return nil, errors.Wrapf(err, "%s create_toAcc %v", s.errCtx, cr)
	}

	if _, err := s.create(*data); err != nil {
		return nil, errors.Wrapf(err, "%s create_repo_create %v", s.errCtx, cr)
	}

	return &Account{}, nil
}
