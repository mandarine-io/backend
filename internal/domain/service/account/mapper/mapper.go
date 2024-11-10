package mapper

import (
	"github.com/mandarine-io/Backend/internal/domain/dto"
	"github.com/mandarine-io/Backend/internal/persistence/model"
)

func MapUserEntityToAccountResponse(userEntity *model.UserEntity) dto.AccountOutput {
	return dto.AccountOutput{
		Username:        userEntity.Username,
		Email:           userEntity.Email,
		IsEnabled:       userEntity.IsEnabled,
		IsEmailVerified: userEntity.IsEmailVerified,
		IsPasswordTemp:  userEntity.IsPasswordTemp,
		IsDeleted:       userEntity.DeletedAt != nil,
	}
}
