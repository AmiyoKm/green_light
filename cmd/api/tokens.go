package main

import (
	"net/http"
	"time"

	"github.com/AmiyoKm/green_light/internal/store"
	"github.com/AmiyoKm/green_light/internal/validator"
)

func (app *application) sendActivationEmail(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email string `json:"email"`
	}

	if err := app.readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	v.Check(validator.Matches(payload.Email, validator.EmailRX), "email", "must be valid email address")
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	
	user, err := app.store.Users.GetByEmail(r.Context(), payload.Email)
	if err != nil {
		switch err {
		case store.ErrorNotFound:
			app.notFoundResponse(w, r)
			return
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}
	token, err := app.store.Tokens.New(r.Context(), user.ID, 3*24*time.Hour, store.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.background(func() {
		data := map[string]any{
			"userID":          user.ID,
			"activationToken": token.Plaintext,
		}

		err := app.mailer.Send(user.Email, "user_welcome.tmpl", data)
		if err != nil {
			app.logger.PrintError(err, nil)
		}
	})

	if err := app.writeJSON(w, http.StatusAccepted, envelope{"message": "email sent"}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
