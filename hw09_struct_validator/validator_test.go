package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require" //nolint: depguard
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
			User{
				ID:     "10985464566",
				Name:   "John",
				Age:    15,
				Email:  "johndoegmail.com",
				Role:   "admin",
				Phones: []string{"79047651231"},
			},
			ValidationErrors{
				ValidationError{
					Field: "Age",
					Err:   ErrValidateMin,
				},
				ValidationError{
					Field: "Email",
					Err:   ErrValidateRegexp,
				},
			},
		},
		{
			User{
				ID:     "986762312312377123983794123123771230000",
				Name:   "James",
				Age:    21,
				Email:  "jamesrodrig12309@gmail.com",
				Role:   "stuff",
				Phones: []string{"79817651231"},
			},
			ValidationErrors{
				ValidationError{
					Field: "ID",
					Err:   ErrValidateLen,
				},
			},
		},
		{
			User{
				ID:     "98123",
				Name:   "Oscar",
				Age:    65,
				Email:  "oscar_guler@gmail.com",
				Role:   "stuff",
				Phones: []string{"79821231233"},
			},
			ValidationErrors{
				ValidationError{
					Field: "Age",
					Err:   ErrValidateMax,
				},
			},
		},
		{
			User{
				ID:     "974",
				Name:   "Boris",
				Age:    43,
				Email:  "boris_03_87@gmail.com",
				Role:   "stuff",
				Phones: []string{"7951563489892321"},
			},
			ValidationErrors{
				ValidationError{
					Field: "Phones",
					Err:   ErrValidateLen,
				},
			},
		},
		{
			User{
				ID:     "109",
				Name:   "Chris",
				Age:    33,
				Email:  "chris_drake078@gmail.com",
				Role:   "other",
				Phones: []string{"79514039089"},
			},
			ValidationErrors{
				ValidationError{
					Field: "Role",
					Err:   ErrValidateIn,
				},
			},
		},
		{
			App{
				Version: "1",
			},
			nil,
		},
		{
			App{
				Version: "123456",
			},
			ValidationErrors{
				ValidationError{
					Field: "Version",
					Err:   ErrValidateLen,
				},
			},
		},
		{
			Response{
				Code: 200,
				Body: "{}",
			},
			nil,
		},
		{
			Response{
				Code: 503,
				Body: "{}",
			},
			ValidationErrors{
				ValidationError{
					Field: "Code",
					Err:   ErrValidateIn,
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			var validError ValidationErrors
			t.Parallel()

			err := Validate(tt.in)

			if errors.As(err, &validError) && err.Error() != "" {
				if valErr, ok := err.(ValidationErrors); ok { //nolint: errorlint
					expectValErr := tt.expectedErr.(ValidationErrors) //nolint: errorlint
					require.Equal(t, len(expectValErr), len(valErr))
					for index, errEntry := range valErr {
						expected := expectValErr[index]
						require.True(t, errors.Is(errEntry.Err, expected.Err))
						require.ErrorContains(t, errEntry.Err, expected.Err.Error())
					}
				}
			}
		})
	}
}
