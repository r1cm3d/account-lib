package acc

import (
	"errors"
	"testing"
)

func TestCreate_Error(t *testing.T) {
	cases := []struct {
		name string
		in   Service
		want error
	}{
		{"toAcc error", serviceWithToAccError, errors.New("service create_toAcc: {}: toAcc error")},
		{"repo create error", serviceWithCreateError, errors.New("service create_repo_create: {}: repo create error")},
		{"ofAcc error", serviceWithOfAccError, errors.New("service create_ofAcc: {}: ofAcc error")},
	}
	cr := CreateRequest{}

	for _, tt := range cases {
		_, got := tt.in.Create(cr)
		if got.Error() != tt.want.Error() {
			t.Errorf("Create_Error(%v) got: %v, want: %v", tt.name, got, tt.want)
		}
	}
}

var (
	serviceWithToAccError = Service{
		errCtx: "service",
		toAcc:  func(createRequest CreateRequest) (*data, error) { return nil, errors.New("toAcc error") },
	}
	serviceWithCreateError = Service{
		errCtx:  "service",
		toAcc:   func(createRequest CreateRequest) (*data, error) { return &data{}, nil },
		creator: mockCreatorErr{},
	}
	serviceWithOfAccError = Service{
		errCtx:  "service",
		toAcc:   func(createRequest CreateRequest) (*data, error) { return &data{}, nil },
		creator: mockCreatorOk{},
		ofAcc:   func(data) (*Account, error) { return nil, errors.New("ofAcc error") },
	}
)

type (
	mockCreatorErr struct{}
	mockCreatorOk  struct{}
)

func (m mockCreatorErr) create(_ data) (*data, error) {
	return nil, errors.New("repo create error")
}

func (m mockCreatorOk) create(_ data) (*data, error) {
	return &data{}, nil
}
