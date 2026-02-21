package withzap

import "go.uber.org/zap"

var zapLog = zap.NewNop()
var zapUserToken = "tok123"

func fZap() {
	// --- Logger methods ---
	zapLog.Info("server ready")
	zapLog.Info("Server started")              // want `log message must start with a lowercase letter`
	zapLog.Info("connection failed!!!")        // want `log message must not contain repeated punctuation`
	zapLog.Info("user token: " + zapUserToken) // want "log message may expose sensitive data" "log message may expose sensitive data"

	// --- SugaredLogger format methods ---
	sugar := zapLog.Sugar()
	sugar.Infof("Starting: %s", "world")        // want `log message must start with a lowercase letter`
	sugar.Infof("user token: %s", zapUserToken) // want "log message may expose sensitive data" "log message may expose sensitive data"

	// --- len(parts)==0: динамический аргумент ---
	dynamicMsg := "hello"
	zapLog.Info(dynamicMsg)
}
