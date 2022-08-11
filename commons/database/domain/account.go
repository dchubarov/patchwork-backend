package domain

type UserAccount struct {
	// Login user login name
	Login string
	// Email email address
	Email string
	// Cn common name
	Cn string
}

var (
	EmptyUserAccount = UserAccount{}
)
