package realworld

import "log/slog"

var rwAddr = ":8080"
var rwJwtToken = "eyJ..."
var rwUserID = "u42"

func rwStartup() {
    slog.Info("starting http server", "addr", rwAddr)
    slog.Info("database connected")
    slog.Info("ready to serve requests")
    slog.Info("serving user " + rwUserID)

    slog.Info("Starting http server") // want `log message must start with a lowercase letter`

    slog.Info("user session: " + rwJwtToken) // want "log message may expose sensitive data"
}

func rwHandleRequest() {
    slog.Info("request received")
    slog.Info("request failed!!!")   // want `log message must not contain repeated punctuation`
    slog.Info("request failed 🔥")   // want `log message must not contain emoji`
}