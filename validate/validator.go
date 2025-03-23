package validate

import (
	"encoding/json"
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

func ValidateJson(data []byte) error {
	// just unmarshal it into a map
	value := map[string]any{}
	return json.Unmarshal(data, &value)
}

func newValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.RegisterValidation("method", httpMethodValidator); err != nil {
		panic(fmt.Sprint("ERROR: failed to register 'http.method' validator", err))
	}
	return v
}
