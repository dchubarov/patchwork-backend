package service

import (
	"net/http"
	"time"
)

// AuthSession contain authentication session data
type AuthSession struct {
	Sid        string    `json:"session"` // Sid the session id
	User       string    `json:"-"`       // User contains user id for the session
	Privileged bool      `json:"-"`       // Privileged specifies whether session has privileged permissions
	Created    time.Time `json:"created"` // Created contains session creation time
	Expires    time.Time `json:"expires"` // Expires session expiration time (Unix)
}

// AuthContext contains authentication data
type AuthContext struct {
	Session *AuthSession // Session contains session info
	User    *AccountUser // User authenticated user
	Token   string       // Token contains authentication token
}

// AuthServiceHeaderCredentials authorization data came from 'Authorization' header
const (
	AuthServiceHeaderCredentials = iota
	//AuthServiceFormCredential
)

var (
	ErrServiceAuthFail            = DefineError("auth.fail", "authentication failed due to internal error", http.StatusInternalServerError)
	ErrServiceAuthInvalidData     = DefineError("auth.invalid", "invalid authentication data supplied", http.StatusUnauthorized)
	ErrServiceAuthNoSession       = DefineError("auth.nosession", "user session not found or expired", http.StatusUnauthorized)
	ErrServiceAuthBadCredentials  = DefineError("auth.credentials", "invalid username or password", http.StatusUnauthorized)
	ErrServiceAuthLoginNotAllowed = DefineError("auth.blocked", "login now allowed for user", http.StatusUnauthorized)
)

// AuthService defines methods of authentication service
type AuthService interface {
	// LoginInternal creates a session for internal user
	LoginInternal(privileged bool) (*AuthContext, error)
	// LoginWithCredentials login user with given credentials
	LoginWithCredentials(authorization string, authorizationType int) (*AuthContext, error)
	// Logout log out user
	Logout(aac *AuthContext) error
}
