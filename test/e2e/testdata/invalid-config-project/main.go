package main

import "log"

// Чистый код — ошибка должна прийти от невалидного конфига, а не от анализа.

func main() {
	log.Print("server started")
	log.Print("connection established")
}
