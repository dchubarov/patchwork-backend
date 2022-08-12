package repos

import "time"

// AuthSession contain authentication session data
type AuthSession struct {
	Sid        string    `json:"session"` // Sid the session id
	User       string    `json:"-"`       // User contains user id for the session
	Privileged bool      `json:"-"`       // Privileged specifies whether session has privileged permissions
	Created    time.Time `json:"created"` // Created contains session creation time
	Expires    time.Time `json:"expires"` // Expires session expiration time (Unix)
}

// AuthRepository defines methods allowing to manage authentication data
type AuthRepository interface {
	// AuthFindSession finds existing session
	AuthFindSession(sid string) (*AuthSession, error)
	// AuthNewSession creates a new session
	AuthNewSession(user *AccountUser) (*AuthSession, error)
}
