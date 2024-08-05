package app

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/salmanshahzad/web-go/internal/utils"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (app *Application) httpLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := responseWriter{ResponseWriter: w}

		defer func() {
			if err := recover(); err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
			}

			req := slog.Group("request",
				"ip", r.RemoteAddr,
				"method", r.Method,
				"path", r.URL.Path,
			)
			res := slog.Group("response",
				"latency", time.Now().Sub(start).String(),
				"status", rw.statusCode,
			)
			logger := slog.Default().With(req, res)

			if rw.statusCode >= http.StatusInternalServerError {
				logger.Error("HTTP", "stack", string(debug.Stack()))
			} else {
				logger.Info("HTTP")
			}
		}()

		next.ServeHTTP(&rw, r)
	})
}

func (app *Application) verifyAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := app.sm.GetInt32(r.Context(), "userId")
		if userId == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		user, err := app.db.GetUser(r.Context(), userId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			utils.InternalServerError(w, r, err)
			return
		}

		if err := app.sm.RenewToken(r.Context()); err != nil {
			utils.InternalServerError(w, r, err)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
