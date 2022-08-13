package services

import (
	"context"
	"os"
	"strconv"
	"twowls.org/patchwork/commons/service"
)

var AuthContextKey = "$aac." + strconv.Itoa(os.Getpid())

// GetAuthFromContext returns authentication context associated with given Context
func GetAuthFromContext(ctx context.Context) *service.AuthContext {
	if raw := ctx.Value(AuthContextKey); raw != nil {
		if aac, ok := raw.(*service.AuthContext); ok {
			return aac
		}
	}
	return nil
}
