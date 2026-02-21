package main

import (
	"log"
	"log/slog"
)

var userToken = "secret-token"

func main() {
	// нарушение 1: заглавная буква
	log.Print("Starting server")

	// нарушение 2: не английский текст
	slog.Info("сервер запущен")

	// нарушение 3: эмодзи
	log.Print("server started 🚀")

	// нарушение 4: повторяющаяся пунктуация
	slog.Info("connection failed!!!")

	// нарушение 5: sensitive данные
	log.Print("auth token: " + userToken)
}
