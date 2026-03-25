package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/kyoO-o/Applicant-skill-assessment-system-using-AI/backend/cmd/web/app"
)

func routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RealIP, middleware.Logger, middleware.Recoverer)
	r.Use(SecureHeaders, app.Session.Enable, Authenticate)

	// if app.Mode == "debug" {
	// 	cr.Get("/api/ws", app.FrontendWS.Handler)
	// } else {
	r.With(RequireAuth).Get("/api/ws", app.FrontendWS.Handler)
	// }

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	r.Route("/pub", func(r chi.Router) {
		r.Get("/login", Login)
		r.Get("/authenticate", BAuthenticate)
	})

	r.With(RequireAuth).Route("/api", func(r chi.Router) {
		r.Get("/me", Me)
		r.Get("/logout", Logout)

		// r.Route("/customers", func(r chi.Router) {
		// 	r.Get("/", getCustomerList)
		// 	r.Post("/add", addCustomer)
		// 	r.Get("/{id}", getCustomer)
		// 	r.Put("/{id}", updateCustomer)
		// 	r.Delete("/{id}", deleteCustomer)
		// })
	})

	return r
}
