package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/justinas/nosurf"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/cmd/web/app"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/common/oapi"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/userman"
)

// Create a NoSurf middleware function which uses a customized CSRF cookie with
// the Secure, Path and HttpOnly flags set.
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})
	return csrfHandler
}

func SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}

func IsAuth(r *http.Request) bool {
	isAuth, ok := r.Context().Value(app.ContextKeyIsAuth).(bool)
	if !ok {
		return false
	}
	return isAuth
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !IsAuth(r) {
			oapi.Forbidden(w)
			return
		}
		w.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		exists := app.Session.Exists(r, "accessToken")
		if !exists || app.Session.GetString(r, "accessToken") == "" {
			app.InfoLog.Println("accessToken not found")
			next.ServeHTTP(w, r)
			return
		}

		email := app.Session.GetString(r, "email")
		if len(email) == 0 {
			next.ServeHTTP(w, r)
			return
		}

		user, err := app.Users.GetWithEmail(email)
		if err != nil {
			if errors.Is(err, userman.ErrNotFound) {
				next.ServeHTTP(w, r)
				return
			}
			oapi.ServerError(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), app.ContextKeyIsAuth, true)
		ctx = context.WithValue(ctx, app.ContextKeyAuthCustomer, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !IsAuth(r) {
			oapi.Forbidden(w)
			return
		}

		user := r.Context().Value(app.ContextKeyAuthCustomer).(*userman.User)

		if user.Role != userman.RoleAdmin {
			oapi.Forbidden(w)
			return
		}

		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}
