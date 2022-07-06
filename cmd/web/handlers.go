package main

import (
	"errors"
	"fmt"

	"net/http"
	"strconv"

	"tutorial-go.com/phonebook/pkg/models"
)

// Обработчик главной странице.
// Меняем сигнатуры обработчика home, чтобы он определялся как метод
// структуры *application.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	s, err := app.numbers.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{
		Numbers: s,
	})
}

// Обработчик для отображения содержимого заметки.
func (app *application) showNumber(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.numbers.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, r, "show.page.tmpl", &templateData{
		Number: s,
	})
}

// Обработчик для создания новой заметки.
func (app *application) createNumber(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)

		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	name := "Настя"
	phone := "89992221234"

	id, err := app.numbers.Insert(name, phone)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/number?id=%d", id), http.StatusSeeOther)

}
