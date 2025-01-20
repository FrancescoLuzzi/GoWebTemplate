package middlewares

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

func Combine(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		var res = next
		for _, m := range middlewares {
			res = m(res)
		}
		return res
	}
}
