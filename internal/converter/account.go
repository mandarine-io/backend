package converter

import (
	"github.com/mandarine-io/backend/internal/persistence/entity"
	"github.com/mandarine-io/backend/pkg/model/v0"
)

func MapUserEntityToAccountOutput(userEntity *entity.User) v0.AccountOutput {
	return v0.AccountOutput{
		Username:        userEntity.Username,
		Email:           userEntity.Email,
		IsEnabled:       userEntity.IsEnabled,
		IsEmailVerified: userEntity.IsEmailVerified,
		IsPasswordTemp:  userEntity.IsPasswordTemp,
		IsDeleted:       userEntity.DeletedAt != nil,
	}
}
