package repo

import (
	"fmt"
)

var (
	// User errors
	ErrDuplicateUser = fmt.Errorf("duplicate user")

	// Master Profile errors
	ErrDuplicateMasterProfile       = fmt.Errorf("duplicate master profile")
	ErrUserForMasterProfileNotExist = fmt.Errorf("user for master profile does not exist")
)
