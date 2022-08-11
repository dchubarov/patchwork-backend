package repos

import (
	"golang.org/x/exp/slices"
)

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

// PasswordMatcher is function that returns true if a password supplied elsewhere matches hashedPassword
type PasswordMatcher func(hashedPassword []byte) bool

// AccountRepository provides methods allowing to access and manage account in database
type AccountRepository interface {
	// AccountFindUser finds user account by login or email
	AccountFindUser(login string, lookupByEmail bool) (*AccountUser, bool)

	// AccountFindLoginUser find AccountUser by login or E-mail, additionally check if user can be logged in, including password check
	AccountFindLoginUser(loginOrEmail string, comparePasswordFn PasswordMatcher) (*AccountUser, bool)
}
