package repos

import (
	"context"
	"twowls.org/patchwork/commons/service"
)

// PasswordMatcher is function that returns true if a password supplied elsewhere matches hashedPassword
type PasswordMatcher func(hashedPassword []byte) bool

// AccountRepository provides methods allowing to access and manage account in database
type AccountRepository interface {
	// AccountFindUser finds user account by login or email
	AccountFindUser(ctx context.Context, login string, lookupByEmail bool) *service.UserAccount
	// AccountFindLoginUser find UserAccount by login or E-mail, additionally check if user can be logged in, including password check
	AccountFindLoginUser(ctx context.Context, loginOrEmail string, comparePasswordFn PasswordMatcher) (*service.UserAccount, bool)
}
