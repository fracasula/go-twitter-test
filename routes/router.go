package routes

import (
	"go-twitter-test/container"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

func NewRouter(c container.Container) *chi.Mux {
	router := chi.NewRouter()

	router.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.DefaultCompress,
		middleware.RedirectSlashes,
		middleware.Recoverer,
		middleware.AllowContentType("application/json"),
		middleware.Timeout(30*time.Second),
		render.SetContentType(render.ContentTypeJSON),
		loggerMiddleware(c.Logger()),
		userMiddleware(),
	)

	router.Route("/v1", func(r chi.Router) {
		r.Mount("/messages", NewMessagesRouter(
			c.MessagesRepository(),
			c.UsersRepository(),
			c.TagsRepository(),
			c.Logger(),
		))
	})

	return router
}
