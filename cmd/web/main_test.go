package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var app *application

func TestMain(m *testing.M) {
	dsn := "web:pass@/test_pb?parseTime=true"
	openDB(dsn)

	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func ensureTableExists() {
	if _, err := app.numbers.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	app.numbers.DB.Exec("DELETE FROM numbers")
	app.numbers.DB.Exec("ALTER SEQUENCE idx_snippets_created RESTART WITH 1")
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS numbers
(
	id int NOT NULL AUTO_INCREMENT,
	name varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
	phone varchar(12) COLLATE utf8mb4_unicode_ci NOT NULL,
	created datetime NOT NULL,
	PRIMARY KEY (id),
	KEY idx_numbers_created (created)
  )`

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	app.routes().ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
