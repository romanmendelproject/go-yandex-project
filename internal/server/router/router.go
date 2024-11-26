package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/romanmendelproject/go-yandex-project/internal/server/config"
	"github.com/romanmendelproject/go-yandex-project/internal/server/handlers"
	"github.com/romanmendelproject/go-yandex-project/internal/server/jwt"
	"github.com/romanmendelproject/go-yandex-project/internal/server/middlewares/auth"
	"github.com/romanmendelproject/go-yandex-project/internal/server/middlewares/compress"
	"github.com/romanmendelproject/go-yandex-project/internal/server/middlewares/hash"
	"github.com/romanmendelproject/go-yandex-project/internal/server/middlewares/logger"
)

func NewRouter(handler *handlers.ServiceHandlers, cfg config.Config, token *jwt.JWT) *chi.Mux {
	r := chi.NewRouter()
	r.Use(logger.RequestLogger)
	r.Use(compress.GzipMiddleware)

	r.Post("/", handlers.HandleBadRequest)

	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", handler.RegisterUser)
		r.Post("/login", handler.LoginUser)

		r.Route("/value", func(r chi.Router) {
			r.Use(auth.IsAuthorized(token))
			r.Route("/cred", func(r chi.Router) {
				r.Get("/{name}", handler.GetCredValue)
				r.Get("/*", handlers.HandleBadRequest)
			})
			r.Route("/text", func(r chi.Router) {
				r.Get("/{name}", handler.GetTextValue)
				r.Get("/*", handlers.HandleBadRequest)
			})
			r.Route("/byte", func(r chi.Router) {
				r.Get("/{name}", handler.GetByteValue)
				r.Get("/*", handlers.HandleBadRequest)
			})
			r.Route("/card", func(r chi.Router) {
				r.Get("/{name}", handler.GetCardValue)
				r.Get("/*", handlers.HandleBadRequest)
			})
		})
		r.Route("/update", func(r chi.Router) {
			r.Use(middleware.AllowContentType("application/json"))

			if cfg.DB.Key != "" {
				r.Use(hash.HashMiddleware(cfg.DB.Key))
			}
			r.Route("/cred", func(r chi.Router) {
				r.Post("/", handler.SetCredValue)
			})
			r.Route("/text", func(r chi.Router) {
				r.Post("/", handler.SetTextValue)
			})
			r.Route("/byte", func(r chi.Router) {
				r.Post("/", handler.SetByteValue)
			})
			r.Route("/card", func(r chi.Router) {
				r.Post("/", handler.SetCardValue)
			})
		})

	})

	r.Get("/ping", handler.Ping)

	return r
}
