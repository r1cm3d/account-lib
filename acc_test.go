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
		{"toAcc error", serviceWithMapError, errors.New("service create_toAcc {}: toAcc error")},
		{"repo create error", serviceWithCreateError, errors.New("service create_repo_create {}: repo create error")},
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
	serviceWithMapError = Service{
		errCtx: "service",
		toAcc:  func(createRequest CreateRequest) (*data, error) { return nil, errors.New("toAcc error") },
	}
	serviceWithCreateError = Service{
		errCtx:  "service",
		toAcc:   func(createRequest CreateRequest) (*data, error) { return &data{}, nil },
		creator: mockCreator{},
	}
)

type mockCreator struct{}

func (m mockCreator) create(_ data) (*data, error) {
	return nil, errors.New("repo create error")
}
