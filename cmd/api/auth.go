package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/high-la/gopher-social/internal/mailer"
	"github.com/high-la/gopher-social/internal/store"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=4,max=16"`
}

type UserWithToken struct {
	*store.User
	Token string `json:"token"`
}

// registerUserHandler godoc
//
//	@Summary		Registers a user
//	@Description	Registers a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload	true	"User credentials"
//
//	@Success		201		{object}	UserWithToken		"User registered"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/user [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {

	var payload RegisterUserPayload

	err := readJSON(w, r, &payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = Validate.Struct(payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
	}
	// hash user password
	err = user.Password.Set(payload.Password)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	// .
	plainToken := uuid.New().String()

	// Store
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	// store user
	err = app.store.Users.CreateAndInvite(ctx, user, hashToken, app.config.mail.expiry)
	if err != nil {

		switch err {
		case store.ErrDuplicateEmail:
			app.badRequestResponse(w, r, err)
		case store.ErrDuplicateUsername:
			app.badRequestResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	userWithToken := UserWithToken{
		User:  user,
		Token: plainToken,
	}

	activationURL := fmt.Sprintf("%s/confirm/%s", app.config.frontendURL, plainToken)

	isProdEnv := app.config.env == "production"
	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.Username,
		ActivationURL: activationURL,
	}

	// send email
	err = app.mailer.Send(mailer.UserWelcomeTemplate, user.Username, user.Email, vars, !isProdEnv)
	if err != nil {
		app.logger.Errorw("error sending welcome email", "error", err)

		// Rollback user creation if sending mail failed
		// Use Background context for the cleanup to ensure it has time to run

		err = app.store.Users.Delete(context.Background(), user.ID)
		if err != nil {
			app.logger.Errorw("error deleting user during user registration roll back", "error", err)
		}

		app.internalServerError(w, r, err)
		return
	}

	err = app.jsonResponse(w, http.StatusCreated, userWithToken)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}
