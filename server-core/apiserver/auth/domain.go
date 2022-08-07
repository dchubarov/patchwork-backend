package auth

type userRole string

type userMembership struct {
	Team       string   `json:"team"`
	CommonName string   `json:"cn,omitempty"`
	Role       userRole `json:"role"`
}

type userEntry struct {
	Login        string           `json:"login"`
	Email        string           `json:"email"`
	CommonName   string           `json:"cn,omitempty"`
	PasswordHash string           `json:"-"`
	MemberOf     []userMembership `json:"memberOf"`
}

type loginResponse struct {
	Session string     `json:"session"`
	Expire  int64      `json:"expires"`
	User    *userEntry `json:"user"`
}
