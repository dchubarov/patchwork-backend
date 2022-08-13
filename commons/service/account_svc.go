package service

import (
	"context"
	"golang.org/x/exp/slices"
	"strings"
)

const (
	UserAccountInternal   = "internal"   // UserAccountInternal indicates an internal user
	UserAccountPrivileged = "privileged" // UserAccountPrivileged indicates a privileged user (administrator)
	UserAccountSuspended  = "suspended"  // UserAccountSuspended indicates an internal user
)

// UserAccount contains user account data
type UserAccount struct {
	Login string   `json:"login"` // Login user login name
	Email string   `json:"email"` // Email email address
	Cn    string   `json:"cn"`    // Cn common name
	Flags []string `json:"flags"` // Flags contains user account flags
}

// Is returns true if specified login matches current user
func (u *UserAccount) Is(login string) bool {
	return u != nil && strings.Compare(u.Login, login) == 0
}

// IsInternal returns true if UserAccountInternal flag is set
func (u *UserAccount) IsInternal() bool {
	return slices.Contains(u.Flags, UserAccountInternal)
}

// IsPrivileged returns true if UserAccountPrivileged flag is set
func (u *UserAccount) IsPrivileged() bool {
	return slices.Contains(u.Flags, UserAccountPrivileged)
}

// IsSuspended returns true if UserAccountSuspended flag is set
func (u *UserAccount) IsSuspended() bool {
	return slices.Contains(u.Flags, UserAccountSuspended)
}

// TeamAccount contains team account information
type TeamAccount struct {
	Team string `json:"team"` // Team the team name
	Cn   string `json:"cn"`   // Cn team common name
}

// AccountService provides methods for managing accounts
type AccountService interface {
	// FindUser get user account by login or email
	FindUser(ctx context.Context, loginOrEmail string, lookupByEmail bool) (*UserAccount, error)
}
