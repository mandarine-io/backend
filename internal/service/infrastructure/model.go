package infrastructure

import (
	"github.com/google/uuid"
)

type AccessTokenClaims struct {
	UserID         uuid.UUID
	Username       string
	Email          string
	Role           string
	IsPasswordTemp bool
	IsEnabled      bool
	IsDeleted      bool
	JTI            string
	Exp            int64
}

type RefreshTokenClaims struct {
	UserID uuid.UUID
	JTI    string
	Exp    int64
}
