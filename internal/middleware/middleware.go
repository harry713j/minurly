package middleware

import "github.com/harry713j/minurly/internal/server"

type Middleware struct {
	Auth *AuthMiddleware
}

func NewMiddleware(server *server.Server) *Middleware {
	return &Middleware{
		Auth: NewAuthMiddleware(server),
	}
}
