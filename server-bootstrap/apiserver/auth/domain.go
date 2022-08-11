package auth

import (
	"twowls.org/patchwork/commons/database/repos"
)

type loginResponse struct {
	Session string             `json:"session"`
	Expire  int64              `json:"expires"`
	User    *repos.AccountUser `json:"user"`
}
