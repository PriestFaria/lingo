package clean

import (
	"log"
	"log/slog"
)

var cleanHost = "localhost"
var cleanPort = "8080"

func startup() {
    log.Print("server started")
    log.Print("database connected")
    log.Printf("listening on %s:%s", cleanHost, cleanPort)

    slog.Info("shutdown complete")
    slog.Info("cache invalidated")
    slog.Info("connected to " + cleanHost)
}