package main

import (
	"fmt"
	"net/http"

	"github.com/AmiyoKm/green_light/internal/data"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Title   string   `json:"title"`
		Year    int32    `json:"year"`
		Runtime int32    `json:"runtime"`
		Genres  []string `json:"genres"`
	}
	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
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
