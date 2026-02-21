package main

import (
	"log"
	"log/slog"
)

// Конфиг: first_letter=false, english=true, emoji=false, security=false.
// Проверяем, что:
//   - заглавная буква НЕ детектируется (first_letter выключен)
//   - не-английский текст ДЕТЕКТИРУЕТСЯ (english включён)
//   - emoji НЕ детектируются (emoji выключен)
//   - sensitive vars НЕ детектируются (security выключен)

var token = "tok"

func main() {
	// first_letter отключён — заглавная буква НЕ должна вызвать ошибку
	log.Print("Starting the server")
	slog.Info("Running graceful shutdown")

	// english включён — должна быть ошибка (кириллица = non-ASCII letters)
	log.Print("сервер запущен")
	slog.Info("соединение установлено")

	// security отключён — токен в переменной НЕ должен вызвать ошибку
	log.Print("request id: " + token)
}
