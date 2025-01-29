package oauth

import "errors"

var (
	ErrUserInfoNotReceived = errors.New("user info not received")
)

type UserInfo struct {
	Username        string
	Email           string
	IsEmailVerified bool
}
