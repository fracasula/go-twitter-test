package routes

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/middleware"
)

const userIDKey = "CTX_USER_ID"

func With(r *http.Request) *ctxReader {
	return &ctxReader{request: r}
}

type ctxReader struct {
	request *http.Request
}

func loggerMiddleware(l *log.Logger) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&middleware.DefaultLogFormatter{
		Logger: l, NoColor: false,
	})
}

func userMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// the userID could be an UUID but to keep things simple I'm using an auto incremented integer
			userID, err := strconv.Atoi(r.Header.Get("X-User-ID"))
			if err == nil {
				ctx = context.WithValue(ctx, userIDKey, int64(userID))
			} else {
				ctx = context.WithValue(ctx, userIDKey, int64(0))
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
