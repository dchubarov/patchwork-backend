package services

import (
	"context"
	"encoding/base64"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"twowls.org/patchwork/commons/database/repos"
	"twowls.org/patchwork/commons/service"
	"twowls.org/patchwork/commons/util/singleton"
	"twowls.org/patchwork/server/bootstrap/database"
	"twowls.org/patchwork/server/bootstrap/logging"
)

const (
	AuthSchemeBasic  = "Basic " // note trailing space
	AuthSchemeBearer = "Bearer "
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

var log = logging.WithComponent("service.auth")

// service.AuthService implementation

func (s *authServiceImpl) LoginInternal(context.Context, bool) (*service.AuthContext, error) {
	// TODO not implemented
	return nil, service.ErrServiceAuthFail
}

func (s *authServiceImpl) LoginWithCredentials(ctx context.Context, authorization string, authorizationType int) (*service.AuthContext, error) {
	if authorizationType == service.AuthServiceHeaderCredentials {
		if strings.HasPrefix(authorization, AuthSchemeBearer) {
			tokenString := authorization[len(AuthSchemeBearer):]
			if sid, err := retrieveSessionFromToken(tokenString); err == nil {
				if session := s.authRepo.AuthFindSession(ctx, sid); session != nil {
					if user := s.accountRepo.AccountFindUser(ctx, session.User, false); user != nil {
						if user.IsInternal() || user.IsSuspended() {
							log.Warnf("Login attempt blocked for user %q with flags %v", user.Login, user.Flags)
							return nil, service.ErrServiceAuthLoginNotAllowed
						}
						return &service.AuthContext{Session: session, User: user}, nil
					}
				} else {
					return nil, service.ErrServiceAuthNoSession
				}
			} else {
				log.Error().Msg("Token validation error")
			}
		} else if strings.HasPrefix(authorization, AuthSchemeBasic) {
			if buf, err := base64.StdEncoding.DecodeString(authorization[len(AuthSchemeBasic):]); err == nil {
				if username, password, ok := strings.Cut(string(buf), ":"); ok {
					user, passwordOk := s.accountRepo.AccountFindLoginUser(ctx, username, func(hashedPassword []byte) bool {
						return passwordMatchesHash(hashedPassword, password)
					})

					if user != nil && passwordOk {
						session := s.authRepo.AuthNewSession(ctx, user)
						if session == nil {
							return nil, service.ErrServiceAuthFail
						}

						token, err := buildToken(session)
						if err != nil {
							return nil, service.ErrServiceAuthFail
						}

						return &service.AuthContext{Session: session, User: user, Token: token}, nil
					}

					return nil, service.ErrServiceAuthBadCredentials
				}
			}
		}
	}

	return nil, service.ErrServiceAuthInvalidData
}

func (s *authServiceImpl) Logout(ctx context.Context) error {
	if aac := GetAuthFromContext(ctx); aac != nil {
		if !s.authRepo.AuthDeleteSession(ctx, aac.Session) {
			log.Errorf("Could not delete session %q", aac.Session.Sid)
			return service.ErrServiceAuthFail
		} else {
			return nil
		}
	}

	return service.ErrServiceAuthInvalidData
}

// private

// TODO must be removed, must not use HMAC either
var builtinHmacSecretKey = []byte("Eij3Hah0uiy8ahgahnah7baghoo6Otho")

func buildToken(session *service.AuthSession) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: &jwt.NumericDate{Time: session.Expires},
		IssuedAt:  &jwt.NumericDate{Time: session.Created},
		ID:        session.Sid,
	})

	if signedToken, err := token.SignedString(builtinHmacSecretKey); err == nil {
		return signedToken, nil
	} else {
		return "", err
	}
}

func retrieveSessionFromToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(tk *jwt.Token) (any, error) {
		if _, ok := tk.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, service.ErrServiceAuthFail
		}
		return builtinHmacSecretKey, nil
	})

	if err == nil && token.Valid {
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if sid, ok := claims["jti"].(string); ok {
				return sid, nil
			}
		}
	}

	return "", service.ErrServiceAuthFail
}

func passwordMatchesHash(hash []byte, password string) bool {
	return bcrypt.CompareHashAndPassword(hash, []byte(password)) == nil
}
