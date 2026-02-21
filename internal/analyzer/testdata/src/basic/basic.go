package basic

import (
	"log"
	"log/slog"
)

var password = "secret123"
var token = "abc"
var msg = "world"
var userToken = "tok123"

func getMessage() string { return "hello" }

func f() {
	// --- log ---
	log.Print("server started")
	log.Print("Starting server")            // want `log message must start with a lowercase letter`
	log.Print("запуск сервера")             // want `log message must be in English`
	log.Print("server started 🚀")          // want `log message must not contain emoji`
	log.Print("connection failed!!!")       // want `log message must not contain repeated punctuation`
	log.Print("user auth: " + password)     // want "log message may expose sensitive data" "log message may expose sensitive data"
	log.Print("token: " + token)            // want "log message may expose sensitive data" "log message may expose sensitive data"

	// --- log format methods ---
	log.Printf("Starting: %s", msg)          // want `log message must start with a lowercase letter`
	log.Printf("user token: %s", userToken)  // want "log message may expose sensitive data" "log message may expose sensitive data"

	// --- slog ---
	slog.Info("server ready")
	slog.Info("Server started")              // want `log message must start with a lowercase letter`
	slog.Info("запуск сервера")              // want `log message must be in English`
	slog.Info("server started 🚀")           // want `log message must not contain emoji`
	slog.Info("connection failed!!!")        // want `log message must not contain repeated punctuation`
	slog.Info("user token: " + userToken)    // want "log message may expose sensitive data" "log message may expose sensitive data"

	// --- len(parts)==0 ---
	log.Print(getMessage())  // динамический аргумент → collectPartsFromExpr вернёт nil
	slog.Info(getMessage())  // аналогично для slog
	_ = slog.With()          // 0 аргументов → len(callExpr.Args)==0 в handleSlog

	// --- non-ADD BinaryExpr: collectPartsFromExpr возвращает nil ---
	var x, y int
	log.Print(x * y)
}
