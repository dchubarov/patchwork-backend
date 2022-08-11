package service

import (
	"twowls.org/patchwork/commons/database/repos"
)

// AuthContext contains authentication data
type AuthContext struct {
	// User authenticated user
	User *repos.AccountUser
}

// AuthService defines methods of authentication service
type AuthService interface {
	Login(authorization string) (*AuthContext, error)
}
