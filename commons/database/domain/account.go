package domain

type UserAccount struct {
	// Exists whether account does exist
	Exists bool
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
