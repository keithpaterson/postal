package validate

import "github.com/go-playground/validator/v10"

// singleton validator
var (
	dataValidator *validator.Validate = newValidator()
)

func ValidateStruct(data any) error {
	return dataValidator.Struct(data)
}
func newValidator() *validator.Validate {
	return validator.New(validator.WithRequiredStructEnabled())
}
