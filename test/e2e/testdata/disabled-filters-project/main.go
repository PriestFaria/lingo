package main

import (
	"log"
	"log/slog"
)

// Все фильтры отключены через .lingo.json — ни одна строка не должна
// вызвать диагностику линтера, несмотря на все виды нарушений.

var userToken = "tok"
var password = "secret"

func main() {
	// нарушение first_letter — отключено
	log.Print("Starting server")
	slog.Info("Running service")

	// нарушение english — отключено
	log.Print("serwer стартовал")
	slog.Info("сервер demarre")

	// нарушение security (variable) — отключено
	log.Print("login: " + userToken)
	log.Print("value: " + password)

	// нарушение repeated punctuation — отключено
	slog.Info("connection failed!!!")
	slog.Info("really bad request 🚀")
}
