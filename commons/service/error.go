package service

import "net/http"

type E struct {
	Code     string `json:"code"` // Code contains error code (e.g. a localisation key)
	Msg      string `json:"msg"`  // Msg contains error message
	HttpCode int    `json:"-"`    // HttpCode contains response http code
}

func (e *E) Error() string {
	return e.Msg
}

func DefineError(code string, msg string, httpCodes ...int) *E {
	httpCode := 500
	if len(httpCodes) > 0 {
		httpCode = httpCodes[0]
	}

	return &E{code, msg, httpCode}
}

var (
	ErrServiceUnspecific     = DefineError("common.fail", "server error", http.StatusInternalServerError)
	ErrServiceLoginRequired  = DefineError("common.unauthorized", "the specified resource requires authorisation", http.StatusUnauthorized)
	ErrServiceNoAccess       = DefineError("common.forbidden", "access to resource forbidden", http.StatusForbidden)
	ErrServiceNoSuchResource = DefineError("common.notfound", "resource not found", http.StatusNotFound)
)
