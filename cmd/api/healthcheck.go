package main

import (
	"fmt"
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Status: available")
	fmt.Fprintf(w, "Environment: %s", app.config.env)
	fmt.Fprintf(w, "Version: %s", version)
}
