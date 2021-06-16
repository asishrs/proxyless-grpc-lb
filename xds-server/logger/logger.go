package logger

import (
	"encoding/json"

	"go.uber.org/zap"
)

// Logger to log via zap
var Logger *zap.Logger

func init() {
	rawJSON := []byte(`{
		"level": "debug",
		"development": true,
		"outputPaths": ["stdout"],
		"errorOutputPaths": ["stderr"],
		"encoding": "json",
		"encoderConfig": {
			"timeKet": "ts",
			"levelKey": "level",
			"nameKey": "logger",
			"callerKey": "caller",
			"messageKey": "msg",
			"stacktraceKey": "stacktrace",
			"lineEnding": "",
			"levelEncoder": "",
			"timeEncoder": "iso8601",
			"durationEncoder": "",
			"callerEncoder": ""
		}
	}`)

	var config zap.Config
	if err := json.Unmarshal(rawJSON, &config); err != nil {
		panic(err)
	}

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	Logger = logger
}
