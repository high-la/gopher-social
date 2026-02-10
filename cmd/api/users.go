package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/high-la/gopher-social/internal/store"
)

type userKey string

const userCtx userKey = "user"

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {

	user := getUserFromCtx(r)

	user, err := app.store.Users.GetByID(r.Context(), user.ID)
	if err != nil {

		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	err = app.jsonResponse(w, http.StatusOK, user)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

// .
func (app *application) usersContextMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//
		idParam := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		ctx := r.Context()

		user, err := app.store.Users.GetByID(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromCtx(r *http.Request) *store.User {

	user, _ := r.Context().Value(userCtx).(*store.User)
	return user
}
