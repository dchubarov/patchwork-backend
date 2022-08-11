package auth

import (
	"twowls.org/patchwork/commons/database/repos"
)

type loginResponse struct {
	Expire int64              `json:"expires"`
	User   *repos.AccountUser `json:"user"`
	Token  string             `json:"token"`
}
