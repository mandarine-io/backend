package mapper

import (
	"mandarine/internal/api/persistence/model"
	"mandarine/internal/api/service/account/dto"
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
