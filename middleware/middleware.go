package middleware

import "qrgen/service/auth"

type middleware struct {
	auth auth.Auth
}

func NewMiddleware(auth auth.Auth) middleware {
	return middleware{
		auth: auth,
	}
}
