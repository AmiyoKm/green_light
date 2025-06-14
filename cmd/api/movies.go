package main

import (
	"fmt"
	"net/http"

	"github.com/AmiyoKm/green_light/internal/data"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new movie")

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
	if err := app.writeJson(w, http.StatusOK, envelop{"movie": movie}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
