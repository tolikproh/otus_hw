package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "123",
				Name:   "name",
				Age:    123,
				Email:  "test",
				Role:   "test",
				Phones: []string{"123"},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "ID",
					Err:   ErrValidLength,
				},
				ValidationError{
					Field: "Age",
					Err:   ErrValidMax,
				},
				ValidationError{
					Field: "Email",
					Err:   ErrValidRegExp,
				},
				ValidationError{
					Field: "Role",
					Err:   ErrValidStrNotIn,
				},
				ValidationError{
					Field: "Phones",
					Err:   ErrValidLength,
				},
			},
		},
		{
			in: User{
				ID:     "58bd82c8-918c-4812-8601-50491dd555e7",
				Name:   "test",
				Age:    32,
				Email:  "t@test.ru",
				Role:   "admin",
				Phones: []string{"89261234567"},
			},
			expectedErr: nil,
		},

		{
			in: App{
				Version: "5",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Version",
					Err:   ErrValidLength,
				},
			},
		},
		{
			in: App{
				Version: "12345",
			},
			expectedErr: nil,
		},

		{
			in: Token{
				Header:    nil,
				Payload:   nil,
				Signature: nil,
			},
			expectedErr: nil,
		},
		{
			in: Token{
				Header:    []byte("application: text/json"),
				Payload:   []byte("status: 200"),
				Signature: []byte("signature"),
			},
			expectedErr: nil,
		},

		{
			in: Response{
				Code: 100,
				Body: "",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Code",
					Err:   ErrValidIntNotIn,
				},
			},
		},
		{
			in: Response{
				Code: 200,
				Body: "",
			},
			expectedErr: nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			require.Equal(t, tt.expectedErr, err)
		})
	}
}
