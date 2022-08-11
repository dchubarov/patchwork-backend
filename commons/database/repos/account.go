package repos

import "twowls.org/patchwork/commons/database/domain"

// AccountRepository provides methods allowing to access and manage account in database
type AccountRepository interface {
	// AccountFindUser returns user account
	AccountFindUser(login string) (domain.UserAccount, bool)
}
