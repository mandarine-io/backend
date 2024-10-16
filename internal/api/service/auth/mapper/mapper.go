package mapper

import (
	"mandarine/internal/api/persistence/model"
	"mandarine/internal/api/service/auth/dto"
	"mandarine/pkg/oauth"
)

func MapRegisterRequestToUserEntity(req dto.RegisterInput) *model.UserEntity {
	return &model.UserEntity{
		Username:        req.Username,
		Email:           req.Email,
		Password:        req.Password,
		IsEnabled:       true,
		IsEmailVerified: true,
		IsPasswordTemp:  false,
	}
}

func MapUserInfoToUserEntity(userInfo oauth.UserInfo) *model.UserEntity {
	return &model.UserEntity{
		Username:        userInfo.Username,
		Email:           userInfo.Email,
		IsEnabled:       true,
		IsEmailVerified: userInfo.IsEmailVerified,
		IsPasswordTemp:  true,
	}
}
