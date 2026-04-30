package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/justinas/nosurf"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/cmd/web/app"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/common/oapi"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/companyman"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/jobman"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/pkg/userman"
	"gorm.io/gorm"
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

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if isAllowedOrigin(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Vary", "Origin")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isAllowedOrigin(origin string) bool {
	if origin == "" {
		return false
	}

	allowedOrigins := []string{
		"http://localhost:3000",
		"http://127.0.0.1:3000",
	}

	for _, allowedOrigin := range allowedOrigins {
		if strings.EqualFold(origin, allowedOrigin) {
			return true
		}
	}

	return false
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

func SetChosenCompany(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !IsAuth(r) {
			oapi.Forbidden(w)
			return
		}

		companyID, err := strconv.Atoi(chi.URLParam(r, "CompanyID"))
		if err != nil || companyID <= 0 {
			oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Invalid company ID"})
			return
		}

		user := r.Context().Value(app.ContextKeyAuthCustomer).(*userman.User)
		if user.CompanyID == nil || int(*user.CompanyID) != companyID {
			oapi.Forbidden(w)
			return
		}

		var company companyman.Company
		if err := app.DB.First(&company, companyID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				oapi.CustomError(w, http.StatusNotFound, map[string]string{"message": "Company not found"})
				return
			}
			oapi.ServerError(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), app.ContextKeyChosenCompany, &company)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func SetChosenJob(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !IsAuth(r) {
			oapi.Forbidden(w)
			return
		}

		jobID, err := strconv.Atoi(chi.URLParam(r, "JobID"))
		if err != nil || jobID <= 0 {
			oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Invalid job ID"})
			return
		}

		company, ok := r.Context().Value(app.ContextKeyChosenCompany).(*companyman.Company)
		if !ok {
			oapi.CustomError(w, http.StatusBadRequest, map[string]string{"message": "Company context is required"})
			return
		}

		user := r.Context().Value(app.ContextKeyAuthCustomer).(*userman.User)
		job, err := app.Jobs.GetForRecruiter(jobID, user.ID)
		if err != nil {
			if errors.Is(err, jobman.ErrNotFound) {
				oapi.CustomError(w, http.StatusNotFound, map[string]string{"message": "Job not found"})
				return
			}
			oapi.ServerError(w, err)
			return
		}

		if job.CompanyID != uint(company.ID) {
			oapi.CustomError(w, http.StatusNotFound, map[string]string{"message": "Job not found"})
			return
		}

		ctx := context.WithValue(r.Context(), app.ContextKeyChosenJob, job)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
