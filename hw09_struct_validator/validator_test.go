package hw09structvalidator

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
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
		Code int `validate:"in:200,404,500"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		title       string
		in          interface{}
		expectedErr error
	}{
		{
			title: "valid User",
			in: User{
				ID:     "40e6215d-b5c6-4896-987c-f30f3678f608",
				Name:   "Иван",
				Age:    25,
				Email:  "example@example.com",
				Role:   "admin",
				Phones: []string{"79987654321", "79987654322", "79987654323"},
			},
			expectedErr: nil,
		},
		{
			title: "invalid User",
			in: User{
				ID:     "1",
				Name:   "Иван",
				Age:    16,
				Email:  "example@example@com",
				Role:   "moderator",
				Phones: []string{"+79987654321", "+79987654322", "+79987654323"},
			},
			expectedErr: ValidationErrors{
				{"ID", ErrValidateStringLen},
				{"Age", ErrValidateIntMin},
				{"Email", ErrValidateStringRegexp},
				{"Role", ErrValidateStringIn},
				{"Phones[0]", ErrValidateStringLen},
				{"Phones[1]", ErrValidateStringLen},
				{"Phones[2]", ErrValidateStringLen},
			},
		},
		{
			title: "valid App",
			in: App{
				Version: "0.1.5",
			},
			expectedErr: nil,
		},
		{
			title: "invalid App",
			in: App{
				Version: "0.12.1",
			},
			expectedErr: ValidationErrors{
				{"Version", ErrValidateStringLen},
			},
		},
		{
			title: "valid empty Token",
			in: Token{
				Header:    []byte{},
				Payload:   []byte{},
				Signature: []byte{},
			},
			expectedErr: nil,
		},
		{
			title: "valid Response",
			in: Response{
				Code: 200,
			},
			expectedErr: nil,
		},
		{
			title: "invalid Response",
			in: Response{
				Code: 301,
			},
			expectedErr: ValidationErrors{
				{"Code", ErrValidateIntIn},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d: %s", i, tt.title), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			if tt.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Equal(t, tt.expectedErr.Error(), err.Error())
			}
		})
	}
}

func TestCheckError(t *testing.T) {
	t.Run("check ErrNotStruct", func(t *testing.T) {
		err := Validate(nil)
		require.Truef(t, errors.Is(err, ErrNotStruct), "actual err - %v", err)

		err = Validate("string value")
		require.Truef(t, errors.Is(err, ErrNotStruct), "actual err - %v", err)
	})

	t.Run("check ErrIncorrectRule", func(t *testing.T) {
		err := Validate(struct {
			Name string `validate:"len:6:25"`
		}{})
		require.Truef(t, errors.Is(err, ErrIncorrectRule), "actual err - %v", err)

		err = Validate(struct {
			Age int `validate:"in:двадцать"`
		}{})
		require.Truef(t, errors.Is(err, ErrIncorrectRule), "actual err - %v", err)

		err = Validate(struct {
			Age int `validate:"len:6"`
		}{})
		require.Truef(t, errors.Is(err, ErrIncorrectRule), "actual err - %v", err)

		err = Validate(struct {
			Age int `validate:"min:двадцать"`
		}{})
		require.Truef(t, errors.Is(err, ErrIncorrectRule), "actual err - %v", err)

		err = Validate(struct {
			Age int `validate:"max:20.1"`
		}{})
		require.Truef(t, errors.Is(err, ErrIncorrectRule), "actual err - %v", err)
	})

	t.Run("check ErrIncorrectType", func(t *testing.T) {
		err := Validate(struct {
			Age float64 `validate:"max:20.1"`
		}{})
		require.Truef(t, errors.Is(err, ErrIncorrectType), "actual err - %v", err)
	})
}
