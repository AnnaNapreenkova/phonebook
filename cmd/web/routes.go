package main

import (
	"net/http"
)

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/number", app.showNumber)
	mux.HandleFunc("/createForm", app.createForm)
	mux.HandleFunc("/number/create", app.createNumber)
	mux.HandleFunc("/delete", app.deleteNumber)
	mux.HandleFunc("/edit", app.editPage)
	mux.HandleFunc("/confirmEdit", app.editNumber)
	mux.HandleFunc("/search", app.searchNumber)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static", http.StripPrefix("/static", fileServer))

	return mux
}
