package main

import (
	"log"
	"log/slog"
)

// Кастомные keywords из .lingo.json: cvv, ssn, otp.

var cvvCode = "1234"     // имя переменной содержит кастомный keyword "cvv"
var ssnNumber = "000"    // имя переменной содержит кастомный keyword "ssn"
var requestID = "req-1"  // безопасное имя — не должно срабатывать

func main() {
	// кастомный keyword "cvv" через имя переменной
	log.Print("card data: " + cvvCode)

	// кастомный keyword "ssn" через имя переменной
	slog.Info("user info: " + ssnNumber)

	// кастомный keyword "otp" в литерале
	log.Print("otp code sent")

	// встроенный keyword "password" по-прежнему работает
	log.Print("password reset initiated")

	// чистые переменные — не должны срабатывать
	log.Print("request: " + requestID)
}
