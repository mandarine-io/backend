package v0

//////////////////// Login //////////////////////

type LoginInput struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type JwtTokensOutput struct {
	AccessToken  string `json:"accessToken" format:"jwt" binding:"required"`
	RefreshToken string `json:"refreshToken" format:"jwt" binding:"required"`
}

//////////////////// Register //////////////////////

type RegisterInput struct {
	Username string `json:"username" format:"username" binding:"required,username"`
	Email    string `json:"email" format:"email" binding:"required,email"`
	Password string `json:"password" format:"zxcvbn" binding:"required,zxcvbn"`
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

//////////////////// Refresh tokens ////////////////////

type RefreshTokensInput struct {
	RefreshToken string `json:"refreshToken" format:"jwt" binding:"required"`
}

//////////////////// Recovery password ////////////////////

type RecoveryPasswordInput struct {
	Email string `json:"email" format:"email" binding:"required,email"`
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

type GetConsentPageURLOutput struct {
	ConsentPageURL string `json:"consentPageURL" format:"uri" binding:"required,uri"`
	OauthState     string `json:"oauthState" binding:"required"`
}

type FetchUserInfoInput struct {
	Code string `json:"code" binding:"required"`
}

type SocialLoginCallbackInput struct {
	Code  string `json:"code" binding:"required"`
	State string `json:"state" binding:"required"`
}
