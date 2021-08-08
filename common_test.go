package account

import (
	"reflect"
	"testing"
)

func assert(t *testing.T, propName string, got, want interface{}) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Entity_%s got: %v, want: %v", propName, got, want)
	}
}
