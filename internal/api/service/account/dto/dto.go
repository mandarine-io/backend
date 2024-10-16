package dto

import "time"

//////////////////// Account ////////////////////

type UpdateUsernameInput struct {
	Username string `json:"username" format:"username" binding:"required,username"`
}

type UpdateEmailInput struct {
	Email string `format:"email" json:"email" binding:"required,email"`
}

type SetPasswordInput struct {
	Password string `json:"password" format:"zxcvbn" binding:"required,zxcvbn"`
}

type UpdatePasswordInput struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" format:"zxcvbn" binding:"required,zxcvbn"`
}

type AccountOutput struct {
	Username        string `json:"username" binding:"required,max=255,min=1"`
	Email           string `json:"email" format:"email" binding:"required,email"`
	IsEnabled       bool   `json:"isEnabled" binding:"required"`
	IsEmailVerified bool   `json:"isEmailVerified" binding:"required"`
	IsPasswordTemp  bool   `json:"isPasswordTemp" binding:"required"`
	IsDeleted       bool   `json:"isDeleted" binding:"required"`
}

//////////////////// Email Verify ////////////////////

type SendEmailParams struct {
	Email string `json:"email" format:"email" binding:"required,email"`
}

type VerifyEmailInput struct {
	OTP   string `json:"otp" binding:"required"`
	Email string `json:"email" format:"email" binding:"required,email"`
}

type EmailVerifyCache struct {
	OTP       string    `json:"otp"`
	Email     string    `json:"email"`
	ExpiredAt time.Time `json:"expiredAt"`
}

type EmailVerifyTemplateArgs struct {
	Email string
	TTL   int
	OTP   string
}
