package vault

import (
	"reflect"
	"testing"
)

func TestNewSecretID(t *testing.T) { //nolint:paralleltest
	type test struct {
		name          string
		argument      string
		id            *secretID
		expectedValid bool
	}

	tests := []test{
		{"AWS format", "secret/hello#foo", &secretID{
			Path:      "secret/hello",
			Key:       "foo",
			awsFormat: true,
		}, true},
		{"invalid AWS format with 2 #", "secret/hello##foo", nil, false},
		{"invalid AWS format ending with #", "secret/hello#", nil, false},
		{"Invalid AWS format ending with slash", "secret/hello#foo/", nil, false},
		{"Invalid AWS format starting with slash", "/secret/hello#foo", nil, false},
		{"invalid AWS format ending with slash after #", "secret/hello#foo/error", nil, false},
		{"invalid AWS format with double slash#", "secret//hello#foo", nil, false},

		{"Keepass format", "secret/hello/foo", &secretID{
			Path:      "secret/hello",
			Key:       "foo",
			awsFormat: false,
		}, true},
		{"Invalid Keepass format ending with slash", "secret/hello/foo/", nil, false},
		{"Invalid Keepass format starting with slash", "/secret/hello/foo", nil, false},
		{"Invalid Keepass format with double slash", "secret//hello/foo", nil, false},

		{"Invalid without slash", "secret", nil, false},
	}

	for _, tt := range tests { //nolint:paralleltest
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			got, err := newSecretID(tt.argument)

			if tt.expectedValid && err != nil {
				t.Errorf("Unexpected error for %v: %v", tt.argument, err)
			} else if !tt.expectedValid && err == nil {
				t.Errorf("Error expected for %v", tt.argument)
			}

			if !reflect.DeepEqual(got, tt.id) {
				t.Errorf("newSecretID() = %v, secretId %v", got, tt.id)
			}
		})
	}
}
