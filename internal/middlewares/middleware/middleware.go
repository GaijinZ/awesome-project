package middleware

import (
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func ChainMiddleware(h http.HandlerFunc, mw ...Middleware) http.HandlerFunc {
	if len(mw) < 1 {
		return h
	}

	wrapped := h

	for i := len(mw) - 1; i >= 0; i-- {
		wrapped = mw[i](wrapped)
	}

	return wrapped
}
