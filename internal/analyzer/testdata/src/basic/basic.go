package basic

import "log"

var password = "secret123"
var token = "abc"

func f() {
	log.Print("server started")
	log.Print("Starting server")            // want `log message must start with a lowercase letter`
	log.Print("запуск сервера")             // want `log message must be in English`
	log.Print("server started 🚀")          // want `log message must not contain emoji`
	log.Print("connection failed!!!")       // want `log message must not contain repeated punctuation`
	log.Print("user auth: " + password)     // want "log message may expose sensitive data" "log message may expose sensitive data"
	log.Print("token: " + token)            // want "log message may expose sensitive data" "log message may expose sensitive data"
}
