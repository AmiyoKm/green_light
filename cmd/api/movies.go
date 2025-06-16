package main

import (
	"fmt"
	"net/http"

	"github.com/AmiyoKm/green_light/internal/data"
	"github.com/AmiyoKm/green_light/internal/validator"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	movie := &data.Movie{
		Title:   payload.Title,
		Genres:  payload.Genres,
		Runtime: payload.Runtime,
		Year:    payload.Year,
	}

	v := validator.New()
	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%v\n", payload)

}
func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	movie := data.Movie{
		ID:      id,
		Title:   "Casablanca",
		Runtime: 102,
		Genres:  []string{"drama", "romance", "war"},
		Version: 1,
	}
	if err := app.writeJSON(w, http.StatusOK, envelop{"movie": movie}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
