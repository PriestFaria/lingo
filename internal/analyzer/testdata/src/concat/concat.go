package concat

import (
	"log"
	"log/slog"
)

var cHost = "localhost"
var cPort = "8080"
var apiSecret = "abc123"
var cUserToken = "tok"

func fConcat() {
    log.Print("connected to " + cHost)
    log.Print("listening on " + cHost + ":" + cPort)
    slog.Info("serving on " + cHost + ":" + cPort)

    log.Print("Connected to " + cHost)  // want `log message must start with a lowercase letter`
    slog.Info("Serving on " + cHost)    // want `log message must start with a lowercase letter`

    log.Print("value: " + apiSecret)  // want "log message may expose sensitive data"
    log.Print("auth: " + cUserToken)  // want "log message may expose sensitive data" "log message may expose sensitive data"

    log.Print("host " + cHost + " is ready")
    log.Print("Host " + cHost + " is ready") // want `log message must start with a lowercase letter`
}