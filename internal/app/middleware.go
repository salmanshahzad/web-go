package app

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/salmanshahzad/web-go/internal/utils"
)

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
