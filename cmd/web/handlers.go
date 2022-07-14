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

	s, err := app.numbers.AllRecords()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{
		Numbers: s,
	})
}

// Обработчик для отображения содержимого
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

// Обработчик поиска
func (app *application) searchNumber(w http.ResponseWriter, r *http.Request) {

	inp := r.FormValue("search")

	s, err := app.numbers.Search(inp, inp)
	if err != nil {
		app.serverError(w, err)
	}

	app.render(w, r, "search.page.tmpl", &templateData{
		Numbers: s,
	})
}

// Отображение формы добавления
func (app *application) createForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		Number: nil,
	})
}

// Обработчик для создания новой записи.
func (app *application) createNumber(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)

		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	name := r.FormValue("name")
	phone := r.FormValue("phone")

	id, err := app.numbers.Insert(name, phone)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/number?id=%d", id), http.StatusSeeOther)
}

// Отображение страницы редактирования
func (app *application) editPage(w http.ResponseWriter, r *http.Request) {
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

	app.render(w, r, "edit.page.tmpl", &templateData{
		Number: s,
	})
}

// Обработчик редактирования
func (app *application) editNumber(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)

		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	newname := r.FormValue("name")
	newphone := r.FormValue("phone")

	err = app.numbers.Edit(newname, newphone, id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/number?id=%d", id), http.StatusSeeOther)
}

// Удаление
func (app *application) deleteNumber(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	err = app.numbers.Delete(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
