package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"tutorial-go.com/phonebook/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
)

// Добавляем поле numbers в структуру application. Это позволит
// сделать объект NumberModel доступным для обработчиков.
type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	numbers       *mysql.NumberModel
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":4000", "Сетевой адрес веб-сервера")
	dsn := flag.String("dsn", "web:pass@/phonebook?parseTime=true", "Название MySQL источника данных")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	// Инициализируем экземпляр mysql.NumberModel и добавляем его в зависимостях.
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		numbers:       &mysql.NumberModel{DB: db},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Запуск сервера на http://127.0.0.1%s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

//TODO
//Ограничение просмотра файлов из директории

//type neuteredFileSystem struct {
//	fs http.FileSystem
//}

//func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
//	f, err := nfs.fs.Open(path)
//	if err != nil {
//		return nil, err
//	}
//
//	s, err := f.Stat()
//	if s.IsDir() {
//		index := filepath.Join(path, "index.html")
//		if _, err := nfs.fs.Open(index); err != nil {
//			closeErr := f.Close()
//			if closeErr != nil {
//				return nil, closeErr
//			}
//
//			return nil, err
//		}
//	}
//
//	return f, nil
//}
