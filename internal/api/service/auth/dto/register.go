package dto

import "time"

type RegisterInput struct {
	Username string `json:"username" format:"username" binding:"required,username"`
	Email    string `json:"email" format:"email" binding:"required,email"`
	Password string `json:"password" format:"zxcvbn" binding:"required,zxcvbn"`
}

type RegisterCache struct {
	User      RegisterInput `json:"user"`
	OTP       string        `json:"otp"`
	ExpiredAt time.Time     `json:"expiredAt"`
}

type RegisterConfirmTemplateArgs struct {
	Email string
	TTL   int
	OTP   string
}

type RegisterConfirmInput struct {
	OTP   string `json:"otp" binding:"required"`
	Email string `json:"email" format:"email" binding:"required,email"`
}
