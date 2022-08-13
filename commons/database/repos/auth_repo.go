package repos

import (
	"context"
	"twowls.org/patchwork/commons/service"
)

// AuthRepository defines methods allowing to manage authentication data
type AuthRepository interface {
	// AuthFindSession finds existing session
	AuthFindSession(ctx context.Context, sid string) *service.AuthSession
	// AuthNewSession creates a new session
	AuthNewSession(ctx context.Context, user *service.UserAccount) *service.AuthSession
	// AuthDeleteSession deletes the specified session
	AuthDeleteSession(ctx context.Context, session *service.AuthSession) bool
}
