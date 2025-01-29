package validator

import (
	"github.com/go-playground/validator/v10"
	"strconv"
	"strings"
)

func PointValidator(fl validator.FieldLevel) bool {
	pointRaw, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	coords := strings.Split(pointRaw, ",")
	if len(coords) != 2 {
		return false
	}

	for _, coord := range coords {
		_, err := strconv.ParseFloat(coord, 64)
		if err != nil {
			return false
		}
	}

	return true
}
