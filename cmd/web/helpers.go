package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("Шаблон %s не существует!", name))
		return
	}

	err := ts.Execute(w, td)
	if err != nil {
		app.serverError(w, err)
	}
}

// Помощник serverError записывает сообщение об ошибке в errorLog и
// затем отправляет пользователю ответ 500 "Внутренняя ошибка сервера
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

}

//TODO
// Помощник clientError отправляет определенный код состояния и соответствующее описание
// пользователю. Нужен чтобы отправлять ответы вроде 400 "Bad Request", когда есть проблема с пользовательским запросом.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// Помощник notFound. Это просто
// удобная оболочка вокруг clientError, которая отправляет пользователю ответ "404 Страница не найдена".
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
