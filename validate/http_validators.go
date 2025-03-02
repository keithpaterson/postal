package validate

import (
	"net/http"

	"github.com/go-playground/validator/v10"
)

func httpMethodValidator(fl validator.FieldLevel) bool {
	switch fl.Field().String() {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodGet, http.MethodDelete:
		return true
	}
	return false
}
