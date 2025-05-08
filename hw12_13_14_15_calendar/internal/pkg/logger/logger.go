package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/pkg/config"
)

type Logger struct {
	config *config.Config
}

type ContextKey string

func ProvideLogger(config *config.Config) *Logger {
	const (
		_ = iota
		levelDebug
		levelInfo
		levelWarn
		levelError
	)
	logLevels := map[string]int{
		"DEBUG": levelDebug,
		"INFO":  levelInfo,
		"WARN":  levelWarn,
		"ERROR": levelError,
	}

	logLevel, ok := logLevels[strings.ToUpper(config.Logger.LogLevel)]
	if !ok {
		logLevel = levelInfo
	}

	if logLevel > levelDebug {
		config.Logger.DisableDebugLog = true
	}
	if logLevel > levelInfo {
		config.Logger.DisableInfoLog = true
	}
	if logLevel > levelWarn {
		config.Logger.DisableWarnLog = true
	}
	if logLevel > levelError {
		config.Logger.DisableErrorLog = true
	}

	return &Logger{
		config: config,
	}
}

func (l *Logger) Error(ctx context.Context, msg string) {
	if !l.config.Logger.DisableErrorLog {
		printError(ctx, "ERROR", msg, l.config.App.Stage)
	}
}

func (l *Logger) Warn(ctx context.Context, msg string) {
	if !l.config.Logger.DisableWarnLog {
		printError(ctx, "WARN", msg, l.config.App.Stage)
	}
}

func (l *Logger) Info(ctx context.Context, msg string) {
	if !l.config.Logger.DisableInfoLog {
		printError(ctx, "INFO", msg, l.config.App.Stage)
	}
}

func (l *Logger) Debug(ctx context.Context, msg string) {
	if !l.config.Logger.DisableDebugLog {
		printError(ctx, "DEBUG", msg, l.config.App.Stage)
	}
}

func (l *Logger) GetContextKey(key string) ContextKey {
	return ContextKey(key)
}

func printError(ctx context.Context, level string, msg string, stage string) {
	timestamp := time.Now().Format(time.RFC3339)

	log := make(map[string]interface{})
	log["timestamp"] = timestamp
	log["level"] = level
	log["stage"] = stage
	log["message"] = msg
	log["context"] = getStringContext(ctx)

	jsonData, _ := json.Marshal(log)
	fmt.Println(string(jsonData))
}

func getStringContext(ctx context.Context) string {
	ctxKeys := []ContextKey{
		"client_ip",
		"http_method",
		"http_path",
		"http_version",
		"http_response_code",
		"latency",
		"user_agent",
	}

	var ctxString string

	if ctx != nil {
		for _, key := range ctxKeys {
			ctxKeyValue, ok := ctx.Value(key).(string)
			if ok {
				ctxString = ctxString + "[" + ctxKeyValue + "]"
			}
		}
	}

	return ctxString
}
