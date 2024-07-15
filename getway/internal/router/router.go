package router

import (
	"getway/internal/models"
	"getway/internal/modules/controller"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	"net/http"
)

type Router struct {
	*chi.Mux
}

func (rout *Router) InitRoutes(c *controller.Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Route("/api", func(r chi.Router) {
		r.Post("/login", c.Auth.HandleLogin)
		r.Post("/register", c.Auth.HandleRegister)
		r.Get("/profile/{id}", c.User.Profile)
		r.Route("/payment", func(r chi.Router) {
			r.Use(jwtauth.Verifier(models.TokenAuth))
			r.Use(jwtauth.Authenticator)
			r.Post("/deposit", c.Pay.HandleDeposit)

		})
	})

	r.Get("/swagger", SwaggerUI)
	r.Get("/docs/*", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.json")
	})

	return r
}
