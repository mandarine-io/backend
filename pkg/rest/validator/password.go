package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/nbutton23/zxcvbn-go"
)

func ZxcvbnPasswordValidator(fl validator.FieldLevel) bool {
	password, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	entropyMatch := zxcvbn.PasswordStrength(password, nil)
	return entropyMatch.Score >= 3
}
