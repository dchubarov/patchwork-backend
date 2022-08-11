package repos

import (
	"golang.org/x/exp/slices"
)

// PasswordMatcher is function that returns true if a password supplied elsewhere matches hashedPassword
type PasswordMatcher func(hashedPassword []byte) bool

const (
	AccountUserInternal   = "internal"   // AccountUserInternal indicates an internal user
	AccountUserPrivileged = "privileged" // AccountUserPrivileged indicates a privileged user (administrator)
	AccountUserSuspended  = "suspended"  // AccountUserSuspended indicates an internal user
)

// AccountUser contains user account data
type AccountUser struct {
	Login string   `json:"login"` // Login user login name
	Email string   `json:"email"` // Email email address
	Cn    string   `json:"cn"`    // Cn common name
	Flags []string `json:"flags"` // Flags contains user account flags
}

// IsInternal returns true if AccountUserInternal flag is set
func (a *AccountUser) IsInternal() bool {
	return slices.Contains(a.Flags, AccountUserInternal)
}

// IsPrivileged returns true if AccountUserPrivileged flag is set
func (a *AccountUser) IsPrivileged() bool {
	return slices.Contains(a.Flags, AccountUserPrivileged)
}

// IsSuspended returns true if AccountUserSuspended flag is set
func (a *AccountUser) IsSuspended() bool {
	return slices.Contains(a.Flags, AccountUserSuspended)
}

// AccountUserRepository provides methods allowing to access and manage account in database
type AccountUserRepository interface {
	// AccountUserFind finds user account by login or email
	AccountUserFind(loginOrEmail string) (*AccountUser, bool)

	// AccountFindLoginUser find AccountUser by login
	AccountFindLoginUser(loginOrEmail string, comparePasswordFn PasswordMatcher) (*AccountUser, bool)
}
