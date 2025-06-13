package main

import (
	"fmt"
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Status available")
	fmt.Fprintf(w, "environment : %s\n", app.config.env)
	fmt.Fprintf(w, "version: %s\n", version)
	w.Write([]byte("Status : ok"))
}
