package services

import (
	"context"
	"twowls.org/patchwork/commons/database/repos"
	"twowls.org/patchwork/commons/service"
	"twowls.org/patchwork/commons/singleton"
	"twowls.org/patchwork/server/bootstrap/database"
)

type accountServiceImpl struct {
	accountRepo repos.AccountRepository
}

var accountService = singleton.Lazy(func() *accountServiceImpl {
	return &accountServiceImpl{
		accountRepo: database.Client().(repos.AccountRepository),
	}
})

func Account() service.AccountService {
	return accountService.Instance()
}

// service.AccountService implementation

func (s *accountServiceImpl) FindUser(ctx context.Context, loginOrEmail string, lookupByEmail bool) (*service.UserAccount, error) {
	if aac := GetAuthFromContext(ctx); aac == nil {
		return nil, service.ErrServiceLoginRequired
	} else if !aac.User.IsPrivileged() && !aac.User.Is(loginOrEmail) {
		return nil, service.ErrServiceNoAccess
	}

	if user := s.accountRepo.AccountFindUser(ctx, loginOrEmail, lookupByEmail); user == nil {
		return nil, service.ErrServiceNoSuchResource
	} else {
		return user, nil
	}
}
