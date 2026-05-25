package main

import (
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/cmd/web/app"
)

func routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RealIP, middleware.Recoverer)
	r.Use(CORS, SecureHeaders, app.Session.Enable, Authenticate)

	// Serve uploaded images (avatars, logos)
	if app.Config.StoragePath != "" {
		storagePath := filepath.Clean(app.Config.StoragePath)
		r.Handle("/storage/*", http.StripPrefix("/storage", http.FileServer(http.Dir(storagePath))))
	}

	// WebSocket is outside the Logger group: after the WS upgrade the connection
	// is hijacked and chi's basicWriter would error trying to write a status code.
	r.With(RequireAuth).Get("/api/ws", app.FrontendWS.Handler)

	// All standard HTTP routes use chi's Logger middleware.
	r.Group(func(r chi.Router) {
		r.Use(middleware.Logger)

		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("pong"))
		})

		r.Route("/pub", func(r chi.Router) {
			r.Get("/login", Login)
			r.Post("/login", PasswordLogin)
			r.Post("/register", Register)
			r.Get("/authenticate", BAuthenticate)
			r.Post("/verify-email", verifyEmailHandler)
			r.Post("/resend-verification", resendVerificationHandler)
			r.Post("/forgot-password", forgotPasswordHandler)
			r.Post("/reset-password", resetPasswordHandler)
		})

		r.With(RequireAuth).Route("/api", func(r chi.Router) {
			r.Get("/me", Me)
			r.Put("/me", UpdateMe)
			r.Post("/me/avatar", UploadAvatar)
			r.Put("/me/email", initiateEmailChangeHandler)
			r.Post("/me/verify-email-change", verifyEmailChangeHandler)
			r.Put("/me/password", changePasswordHandler)
			r.Get("/logout", Logout)
			r.Get("/company", GetCompany)
			r.Put("/company", SaveCompany)
			r.Post("/company/logo", UploadCompanyLogo)

			// Job routes
			r.Route("/jobs", func(r chi.Router) {
				r.Get("/", getJobs)
				r.Post("/", saveJob)
				r.With(SetChosenCompany).Route("/{CompanyID}", func(r chi.Router) {
					r.Get("/", getCompanyJobs)
					r.With(SetChosenJob).Route("/{JobID}", func(r chi.Router) {
						r.Get("/", getJob)
						r.Put("/", saveCompanyJob)
						r.Delete("/", deleteJob)
					})
				})
			})

			// CV application routes
			r.Post("/jobs/{JobID}/apply", applyToJob)
			r.Post("/jobs/{JobID}/apply-from-profile", applyFromProfile)
			r.Post("/jobs/{JobID}/analyze", analyzeCV)
			r.Post("/jobs/{JobID}/analyze-from-profile", analyzeFromProfile)
			r.Get("/jobs/{JobID}/applications", listJobApplications)
			r.Get("/applications", listMyApplications)
			r.Get("/applications/{id}", getApplication)
			r.Put("/applications/{id}/status", updateApplicationStatus)
			r.Put("/applications/{id}/interview", scheduleInterview)
			r.Get("/interviews", listInterviewsHandler)

			// CV profile routes
			r.Get("/cv-profile", getCVProfile)
			r.Put("/cv-profile", saveCVProfile)
			r.Post("/cv-profile/parse", parseCVProfile)

			// Google Calendar integration
			r.Get("/integrations/google-calendar/connect", gcalConnect)
			r.Get("/integrations/google-calendar/callback", gcalCallback)
			r.Delete("/integrations/google-calendar/disconnect", gcalDisconnect)
			r.Get("/integrations/google-calendar/status", gcalStatus)

			// Chatbot
			r.Post("/chat", chat)

			// Task routes
			r.Get("/tasks", listTasks)
			r.Post("/tasks", createTask)
			r.Post("/tasks/generate", generateTask)
			r.Get("/tasks/submissions", listAllSubmissions)
			r.Get("/tasks/{id}", getTask)
			r.Put("/tasks/{id}", updateTask)
			r.Delete("/tasks/{id}", deleteTask)
			r.Put("/tasks/{id}/send", sendTask)
			r.Post("/tasks/{id}/submit", submitTask)
			r.Put("/tasks/submissions/{id}/grade", gradeSubmission)
			// r.Post("/tasks/submissions/{id}/ai-grade", aiGradeSubmission)
			r.Get("/tasks/submissions/{id}/file", serveSubmissionFile)
		})
	})

	return r
}
