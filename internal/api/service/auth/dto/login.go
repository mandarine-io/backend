package dto

type LoginInput struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type JwtTokensOutput struct {
	AccessToken  string `json:"accessToken" format:"jwt" binding:"required"`
	RefreshToken string `json:"-"`
}
