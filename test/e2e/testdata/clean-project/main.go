package main

import (
	"log"
	"log/slog"
)

var host = "localhost"
var port = "8080"

func main() {
	log.Print("server starting")
	log.Printf("listening on %s:%s", host, port)
	log.Print("database connected")

	slog.Info("ready to serve requests")
	slog.Info("connected to " + host)
	slog.Info("config loaded")
}
