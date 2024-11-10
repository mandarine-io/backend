package dto

import "time"

//////////////////// Login //////////////////////

type LoginInput struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type JwtTokensOutput struct {
	AccessToken  string `json:"accessToken" format:"jwt" binding:"required"`
	RefreshToken string `json:"-"`
}

//////////////////// Register //////////////////////

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

//////////////////// Social login ////////////////////

type GetConsentPageUrlOutput struct {
	ConsentPageUrl string `json:"consentPageUrl" format:"uri" binding:"required,uri"`
	OauthState     string `json:"oauthState" binding:"required"`
}

type FetchUserInfoInput struct {
	Code string `json:"code" binding:"required"`
}

type SocialLoginCallbackInput struct {
	Code  string `json:"code" binding:"required"`
	State string `json:"state" binding:"required"`
}
