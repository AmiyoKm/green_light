package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/AmiyoKm/green_light/internal/store"
	"github.com/AmiyoKm/green_light/internal/validator"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Title   string        `json:"title"`
		Year    int32         `json:"year"`
		Runtime store.Runtime `json:"runtime"`
		Genres  []string      `json:"genres"`
	}

	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	movie := &store.Movie{
		Title:   payload.Title,
		Runtime: payload.Runtime,
		Year:    payload.Year,
		Genres:  payload.Genres,
	}

	v := validator.New()
	if store.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.store.Movies.Create(r.Context(), movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))
	if err := app.writeJSON(w, http.StatusCreated, envelop{"movie": movie}, headers); err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	movie, err := app.store.Movies.Get(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			app.notFoundResponse(w, r)
			return
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}
	if err := app.writeJSON(w, http.StatusOK, envelop{"movie": movie}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie, err := app.store.Movies.Get(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			app.notFoundResponse(w, r)
			return
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}
	if r.Header.Get("X-Expected-Version") != "" {
		if strconv.FormatInt(int64(movie.Version), 32) != r.Header.Get("X-Expected-Version") {
			app.editConflictResponse(w, r)
			return
		}
	}
	var payload struct {
		Title   *string        `json:"title"`
		Year    *int32         `json:"year"`
		Runtime *store.Runtime `json:"runtime"`
		Genres  []string       `json:"genres"`
	}
	err = app.readJSON(w, r, &payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if payload.Title != nil {
		movie.Title = *payload.Title
	}
	if payload.Year != nil {
		movie.Year = *payload.Year
	}
	if payload.Runtime != nil {
		movie.Runtime = *payload.Runtime
	}
	if payload.Genres != nil {
		movie.Genres = payload.Genres
	}

	v := validator.New()
	if store.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	ctx := r.Context()
	err = app.store.Movies.Update(ctx, movie)

	if err != nil {
		switch err {
		case store.EditConflict:
			app.editConflictResponse(w, r)
			return
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}
	if err := app.writeJSON(w, http.StatusOK, envelop{"movie": movie}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.store.Movies.Delete(r.Context(), id)
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
	err = app.writeJSON(w, http.StatusOK, envelop{"message": "movie successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
