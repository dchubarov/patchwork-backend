package repos

import (
	"twowls.org/patchwork/commons/service"
)

// AuthRepository defines methods allowing to manage authentication data
type AuthRepository interface {
	// AuthFindSession finds existing session
	AuthFindSession(sid string) *service.AuthSession
	// AuthNewSession creates a new session
	AuthNewSession(user *service.AccountUser) *service.AuthSession
	// AuthDeleteSession deletes the specified session
	AuthDeleteSession(session *service.AuthSession) bool
}
