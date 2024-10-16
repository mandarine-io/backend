package random

import (
	"crypto/rand"
	"errors"
	"math/big"
)

var (
	ErrInvalidOtpLength = errors.New("OTP length is negative")
)

func GenerateRandomNumber(length int) (string, error) {
	if length < 0 {
		return "", ErrInvalidOtpLength
	}

	const digits = "0123456789"
	result := make([]byte, length)

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		if num == nil {
			num = big.NewInt(0)
		}

		result[i] = digits[num.Int64()]
	}

	return string(result), nil
}
