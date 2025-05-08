package logger

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/pkg/config"
)

func TestLogger(t *testing.T) {
	conf := &config.Config{
		Logger: config.LoggerConf{
			LogLevel:        "DEBUG",
			DisableErrorLog: false,
			DisableWarnLog:  false,
			DisableInfoLog:  false,
			DisableDebugLog: false,
		},
	}
	logg := ProvideLogger(conf)
	ctx := context.WithValue(context.Background(), logg.GetContextKey("client_ip"), "127.0.0.1")

	tests := []struct {
		level  string
		method func(*Logger, context.Context, string)
	}{
		{"DEBUG", (*Logger).Debug},
		{"INFO", (*Logger).Info},
		{"WARN", (*Logger).Warn},
		{"ERROR", (*Logger).Error},
	}

	for _, testItem := range tests {
		conf.Logger.LogLevel = testItem.level

		t.Run(testItem.level+" log enabled", func(t *testing.T) {
			logg = ProvideLogger(conf)

			message := "Test log " + testItem.level
			output := getStdoutByFunc(func() {
				testItem.method(logg, ctx, message)
			})

			assert.Contains(t, output, `"level":"`+testItem.level+`"`)
			assert.Contains(t, output, `"message":"`+message+`"`)
			assert.Contains(t, output, `[127.0.0.1]`)
		})

		switch testItem.level {
		case "DEBUG":
			conf.Logger.DisableDebugLog = true
		case "INFO":
			conf.Logger.DisableInfoLog = true
		case "WARN":
			conf.Logger.DisableWarnLog = true
		case "ERROR":
			conf.Logger.DisableErrorLog = true
		default:
			t.Fatalf("Unknown log level: %s", testItem.level)
		}

		t.Run(testItem.level+" log disabled", func(t *testing.T) {
			logg = ProvideLogger(conf)

			message := "Test log " + testItem.level
			output := getStdoutByFunc(func() {
				testItem.method(logg, ctx, message)
			})

			require.Equal(t, "", output)
		})
	}
}

func getStdoutByFunc(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}
