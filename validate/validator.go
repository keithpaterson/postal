package validate

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// singleton validator
var (
	dataValidator *validator.Validate = newValidator()
)

func ValidateStruct(data any) error {
	return dataValidator.Struct(data)
}

func newValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.RegisterValidation("method", httpMethodValidator); err != nil {
		panic(fmt.Sprint("ERROR: failed to register 'http.method' validator", err))
	}
	return v
}
