package validator

import (
	"github.com/go-playground/validator/v10"
	"time"
)

func PastDateValidator(fl validator.FieldLevel) bool {
	dateStr, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	date, err := time.Parse(time.DateOnly, dateStr)
	if err != nil {
		return false
	}

	today := time.Now()
	return today.After(date)
}
