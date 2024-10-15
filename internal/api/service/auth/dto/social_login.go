package dto

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
