package main

import (
	"errors"
	"net/http"

	"github.com/AmiyoKm/green_light/internal/store"
	"github.com/AmiyoKm/green_light/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &store.User{
		Name:      payload.Name,
		Email:     payload.Email,
		Activated: false,
	}
	err = user.Password.Set(payload.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	if store.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.store.Users.Create(r.Context(), user)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
			return
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	app.background(func() {
		err = app.mailer.Send(user.Email, "user_welcome.tmpl", user)
		if err != nil {
			app.logger.PrintError(err, nil)
		}
	})

	err = app.writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
