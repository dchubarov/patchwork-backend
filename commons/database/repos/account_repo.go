package repos

import (
	"twowls.org/patchwork/commons/service"
)

// PasswordMatcher is function that returns true if a password supplied elsewhere matches hashedPassword
type PasswordMatcher func(hashedPassword []byte) bool

// AccountRepository provides methods allowing to access and manage account in database
type AccountRepository interface {
	// AccountFindUser finds user account by login or email
	AccountFindUser(login string, lookupByEmail bool) *service.AccountUser
	// AccountFindLoginUser find AccountUser by login or E-mail, additionally check if user can be logged in, including password check
	AccountFindLoginUser(loginOrEmail string, comparePasswordFn PasswordMatcher) (*service.AccountUser, bool)
}
