package middleware

import (
	"context"
	"net/http"
	"social/app"
	"social/internal/core/service"
	"strings"
)

func Auth(ts *service.TokenService, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := &app.Context{
			ResponseWriter: w,
			Request:        r,
		}
		if id, ok := isAuth(ts, ctx.Request); ok {
			_ctx := context.WithValue(ctx.Request.Context(), "userId", id)
			next(ctx.ResponseWriter, ctx.Request.WithContext(_ctx))
		} else {
			ctx.HandleError("please sign in first", http.StatusUnauthorized)
		}
	}
}

func isAuth(ts *service.TokenService, r *http.Request) (int, bool) {
	var token string
	auth := r.Header.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		token = strings.TrimPrefix(auth, "Bearer ")
	}
	id, err := ts.VerifyToken(r.Context(), token)
	return id, err == nil
}
