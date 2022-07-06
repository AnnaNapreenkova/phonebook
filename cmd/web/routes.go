package main

import (
	"net/http"
)

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/number", app.showNumber)
	mux.HandleFunc("/number/create", app.createNumber)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static", http.StripPrefix("/static", fileServer))

	return mux
}
