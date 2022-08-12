package service

import (
	"twowls.org/patchwork/commons/database/repos"
)

// AuthContext contains authentication data
type AuthContext struct {
	Session *repos.AuthSession // Session contains session info
	User    *repos.AccountUser // User authenticated user
	Token   string             // Token contains authentication token
}

const (
	AuthServiceHeaderCredentials = iota // AuthServiceHeaderCredentials authorization data came from 'Authorization' header
	//AuthServiceFormCredential
)

// AuthService defines methods of authentication service
type AuthService interface {
	// LoginInternal creates a session for internal user
	LoginInternal(privileged bool) (*AuthContext, error)

	// LoginWithCredentials login user with given credentials
	LoginWithCredentials(authorization string, authorizationType int) (*AuthContext, error)
}
