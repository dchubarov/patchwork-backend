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
	AuthSchemeBasic  = "Basic " // note trailing space
	AuthSchemeBearer = "Bearer "
)

var (
	ErrAuthInvalidData    = errors.New("invalid authorization data supplied")
	ErrAuthBadCredentials = errors.New("invalid username or password")
	ErrAuthCreateToken    = errors.New("token creation error")
	ErrAuthInvalidToken   = errors.New("invalid token")
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
	if authorizationType == service.AuthServiceHeaderCredentials {
		if strings.HasPrefix(authorization, AuthSchemeBearer) {
			tokenString := authorization[len(AuthSchemeBearer):]
			if sid, err := validateToken(tokenString); err == nil {
				if session, err := s.authRepo.AuthFindSession(sid); err == nil {
					// TODO check if user suspended or internal
					if user, found := s.accountRepo.AccountFindUser(session.User, false); found {
						return &service.AuthContext{Session: session, User: user, Token: tokenString}, nil
					}
				}
			}
		} else if strings.HasPrefix(authorization, AuthSchemeBasic) {
			if buf, err := base64.StdEncoding.DecodeString(authorization[len(AuthSchemeBasic):]); err == nil {
				if username, password, ok := strings.Cut(string(buf), ":"); ok {
					passwordMatcher := func(hashedPassword []byte) bool {
						return passwordMatchesHash(hashedPassword, password)
					}

					if user, found := s.accountRepo.AccountFindLoginUser(username, passwordMatcher); found {
						if session, err := s.authRepo.AuthNewSession(user); err == nil {
							if token, err := buildToken(session); err == nil {
								return &service.AuthContext{Session: session, User: user, Token: token}, nil
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
	}
	return nil, ErrAuthInvalidData
}

// private

var hmacSecret = []byte("Eij3Hah0uiy8ahgahnah7baghoo6Otho")

func buildToken(session *repos.AuthSession) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: &jwt.NumericDate{Time: session.Expires},
		IssuedAt:  &jwt.NumericDate{Time: session.Created},
		ID:        session.Sid,
	})

	if signedToken, err := token.SignedString(hmacSecret); err == nil {
		return signedToken, nil
	} else {
		return "", err
	}
}

func validateToken(tokenString string) (string, error) {
	if token, err := jwt.Parse(tokenString, func(tk *jwt.Token) (any, error) {
		if _, ok := tk.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrAuthInvalidToken
		}
		return hmacSecret, nil
	}); err == nil && token.Valid {
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if sid, ok := claims["jti"].(string); ok {
				return sid, nil
			}
		}
	}

	return "", ErrAuthInvalidToken
}

func passwordMatchesHash(hash []byte, password string) bool {
	return bcrypt.CompareHashAndPassword(hash, []byte(password)) == nil
}
