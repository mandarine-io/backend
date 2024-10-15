package dto

import "time"

//////////////////// Recovery password ////////////////////

type RecoveryPasswordInput struct {
	Email string `json:"email" format:"email" binding:"required,email"`
}

type RecoveryPasswordCache struct {
	OTP       string    `json:"otp"`
	Email     string    `json:"email"`
	ExpiredAt time.Time `json:"expiredAt"`
}

type RecoveryPasswordTemplateArgs struct {
	Email string
	TTL   int
	OTP   string
}

//////////////////// Verify recovery password ////////////////////

type VerifyRecoveryCodeInput struct {
	OTP   string `json:"otp" binding:"required"`
	Email string `json:"email" format:"email" binding:"required,email"`
}

//////////////////// Reset password ////////////////////

type ResetPasswordInput struct {
	OTP      string `json:"otp" binding:"required"`
	Email    string `json:"email" binding:"required,email" format:"email"`
	Password string `json:"password" binding:"required,zxcvbn" format:"zxcvbn"`
}
