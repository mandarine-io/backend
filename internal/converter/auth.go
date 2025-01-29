package converter

import (
	"github.com/mandarine-io/backend/internal/persistence/entity"
	"github.com/mandarine-io/backend/pkg/model/v0"
	"github.com/mandarine-io/backend/third_party/oauth"
)

func MapRegisterInputToUserEntity(input v0.RegisterInput) *entity.User {
	return &entity.User{
		Username:        input.Username,
		Email:           input.Email,
		Password:        input.Password,
		IsEnabled:       true,
		IsEmailVerified: true,
		IsPasswordTemp:  false,
	}
}

func MapUserInfoToUserEntity(userInfo oauth.UserInfo) *entity.User {
	return &entity.User{
		Username:        userInfo.Username,
		Email:           userInfo.Email,
		IsEnabled:       true,
		IsEmailVerified: userInfo.IsEmailVerified,
		IsPasswordTemp:  true,
	}
}
