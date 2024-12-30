package validator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

func UsernameValidator(fl validator.FieldLevel) bool {
	username, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	matched, err := regexp.MatchString("^[a-z][a-z0-9_]{1,255}$", username)
	return err == nil && matched
}
