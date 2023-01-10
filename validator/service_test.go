package validator_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/tanyudii/core-go/errutil"
	"github.com/tanyudii/core-go/validator"
	"testing"
)

type User struct {
	Name      string     `validate:"required"`
	Email     string     `validate:"required,email"`
	Addresses []*Address `validate:"required,dive"`
	Bank      *Bank      `validate:"required"`
}

type Address struct {
	City       string `validate:"required"`
	PostalCode string `validate:"required" name:"Postal code"`
}

type Bank struct {
	Number string `validate:"required"`
	Name   string `validate:"required"`
}

func TestNormalValidate(t *testing.T) {
	v := validator.NewValidator()

	testCases := map[string]struct {
		Request *User
		Error   error
	}{
		"[S] Normal": {
			Request: &User{
				Name:  "John Doe",
				Email: "admin@example.com",
				Addresses: []*Address{
					{City: "Jakarta", PostalCode: "15100"},
					{City: "Tangerang", PostalCode: "15110"},
				},
				Bank: &Bank{
					Number: "9988776655",
					Name:   "Example",
				},
			},
			Error: nil,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			err := v.ValidateStruct(testCase)
			assert.Equal(t, testCase.Error, err)
		})
	}
}

func TestErrorValidate(t *testing.T) {
	v := validator.NewValidator()

	testCases := map[string]struct {
		Request *User
		Error   error
	}{

		"[F] Error Basic Field": {
			Request: &User{
				Name:  "",
				Email: "",
				Addresses: []*Address{
					{City: "Jakarta", PostalCode: "15100"},
					{City: "Tangerang", PostalCode: "15110"},
				},
				Bank: &Bank{
					Number: "9988776655",
					Name:   "Example",
				},
			},
			Error: errutil.NewBadRequestErrorUsingFieldsOrNil(errutil.ErrorField{
				"name":  "Name is a required field",
				"email": "Email is a required field",
			}),
		},
		"[F] Error Nesting N Field": {
			Request: &User{
				Name:  "John Doe",
				Email: "admin@example.com",
				Addresses: []*Address{
					{},
					{},
				},
				Bank: &Bank{
					Number: "9988776655",
					Name:   "Example",
				},
			},
			Error: errutil.NewBadRequestErrorUsingFieldsOrNil(errutil.ErrorField{
				"addresses.0.city":       "City is a required field",
				"addresses.0.postalCode": "Postal code is a required field",
				"addresses.1.city":       "City is a required field",
				"addresses.1.postalCode": "Postal code is a required field",
			}),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			err := v.ValidateStruct(testCase)
			assert.NotNil(t, testCase.Error)
			assert.NotNil(t, err)

			testCaseErrBR := testCase.Error.(*errutil.BadRequestError)
			errBR := err.(*errutil.BadRequestError)
			assert.Equal(t, testCaseErrBR.GetFields(), errBR.GetFields())
		})
	}
}
