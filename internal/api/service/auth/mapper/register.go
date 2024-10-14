package mapper

import (
	"mandarine/internal/api/persistence/model"
	"mandarine/internal/api/service/auth/dto"
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
