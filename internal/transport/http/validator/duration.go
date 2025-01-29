package validator

import (
	"github.com/go-playground/validator/v10"
	"time"
)

func DurationValidator(fl validator.FieldLevel) bool {
	durationStr, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	_, err := time.ParseDuration(durationStr)
	return err == nil
}
