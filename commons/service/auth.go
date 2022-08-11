package service

import (
	"twowls.org/patchwork/commons/database/repos"
)

// AuthContext contains authentication data
type AuthContext struct {
	Session *repos.AuthSession // Session contains session info
	User    *repos.AccountUser // User authenticated user
}

// AuthService defines methods of authentication service
type AuthService interface {
	Login(authorization string) (*AuthContext, error)
}
