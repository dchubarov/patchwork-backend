package services

import (
	"encoding/base64"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"twowls.org/patchwork/commons/database/repos"
	"twowls.org/patchwork/commons/service"
	"twowls.org/patchwork/commons/singleton"
	"twowls.org/patchwork/server/bootstrap/database"
)

const (
	authSchemeBasic = "Basic " // note trailing space
)

var (
	ErrAuthInvalidData    = errors.New("invalid authorization data supplied")
	ErrAuthBadCredentials = errors.New("invalid username or password")
	ErrAuthCreateToken    = errors.New("token creation error")
)

type authServiceImpl struct {
	accountRepo repos.AccountRepository
	authRepo    repos.AuthRepository
}

var authService = singleton.Lazy(func() *authServiceImpl {
	return &authServiceImpl{
		database.Client().(repos.AccountRepository),
		database.Client().(repos.AuthRepository),
	}
})

// Auth returns authorization service instance
func Auth() service.AuthService {
	return authService.Instance()
}

// service.AuthService implementation

func (s *authServiceImpl) LoginInternal(bool) (*service.AuthContext, error) {
	// TODO not implemented
	return nil, nil
}

func (s *authServiceImpl) LoginWithCredentials(authorization string, authorizationType int) (*service.AuthContext, error) {
	if authorizationType == service.AuthServiceHeaderCredentials && strings.HasPrefix(authorization, authSchemeBasic) {
		if buf, err := base64.StdEncoding.DecodeString(authorization[len(authSchemeBasic):]); err == nil {
			if username, password, ok := strings.Cut(string(buf), ":"); ok {
				passwordMatcher := func(hashedPassword []byte) bool {
					return passwordMatchesHash(hashedPassword, password)
				}

				if user, found := s.accountRepo.AccountFindLoginUser(username, passwordMatcher); found {
					if session, err := s.authRepo.AuthNewSession(user); err == nil {
						token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
							ExpiresAt: &jwt.NumericDate{Time: session.Expires},
							IssuedAt:  &jwt.NumericDate{Time: session.Created},
							ID:        session.Sid,
						})

						// TODO testing only: must not use HMAC, must not use static key either
						if signedToken, err := token.SignedString([]byte("Eij3Hah0uiy8ahgahnah7baghoo6Otho")); err == nil {
							return &service.AuthContext{Session: session, User: user, Token: signedToken}, nil
						} else {
							return nil, ErrAuthCreateToken
						}
					} else {
						return nil, err
					}
				} else {
					return nil, ErrAuthBadCredentials
				}
			}
		}
	}

	return nil, ErrAuthInvalidData
}

// private

func passwordMatchesHash(hash []byte, password string) bool {
	return bcrypt.CompareHashAndPassword(hash, []byte(password)) == nil
}
