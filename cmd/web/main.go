package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Создаем структуру `application` для хранения зависимостей всего веб-приложения.
// Пока, что мы добавим поля только для двух логгеров, но
// мы будем расширять данную структуру по мере усложнения приложения.

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	// Создаем новый флаг командной строки, значение по умолчанию: ":4000".
	// Добавляем небольшую справку, объясняющая, что содержит данный флаг.
	// Значение флага будет сохранено в переменной addr.
	addr := flag.String("addr", ":4000", "Сетевой адрес HTTP")

	// Мы вызываем функцию flag.Parse() для извлечения флага из командной строки.
	// Она считывает значение флага из командной строки и присваивает его содержимое
	// переменной. Вам нужно вызвать ее *до* использования переменной addr
	// иначе она всегда будет содержать значение по умолчанию ":4000".
	// Если есть ошибки во время извлечения данных - приложение будет остановлено.
	flag.Parse()

	//логирование
	// Используйте log.New() для создания логгера для записи информационных сообщений. Для этого нужно
	// три параметра: место назначения для записи логов (os.Stdout), строка
	// с префиксом сообщения (INFO или ERROR) и флаги, указывающие, какая
	// дополнительная информация будет добавлена. Обратите внимание, что флаги
	// соединяются с помощью оператора OR |.

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Создаем логгер для записи сообщений об ошибках таким же образом, но используем stderr как
	// место для записи и используем флаг log.Lshortfile для включения в лог
	// названия файла и номера строки где обнаружилась ошибка.

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Инициализируем новую структуру с зависимостями приложения.
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	// Инициализируем FileServer, он будет обрабатывать
	// HTTP-запросы к статическим файлам из папки "./ui/static".
	// Обратите внимание, что переданный в функцию http.Dir путь
	// является относительным корневой папке проекта
	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./static")})
	mux.Handle("/static", http.NotFoundHandler())

	// Используем функцию mux.Handle() для регистрации обработчика для
	// всех запросов, которые начинаются с "/static/". Мы убираем
	// префикс "/static" перед тем как запрос достигнет http.FileServer
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Инициализируем новую структуру http.Server. Мы устанавливаем поля Addr и Handler, так
	// что сервер использует тот же сетевой адрес и маршруты, что и раньше, и назначаем
	// поле ErrorLog, чтобы сервер использовал наш логгер
	// при возникновении проблем.

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	// Значение, возвращаемое функцией flag.String(), является указателем на значение
	// из флага, а не самим значением. Нам нужно убрать ссылку на указатель
	// то есть перед использованием добавьте к нему префикс *. Обратите внимание, что мы используем
	// функцию log.Printf() для записи логов в журнал работы нашего приложения.

	infoLog.Printf("Запуск сервера на %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}
