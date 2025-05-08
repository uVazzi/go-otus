package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		config := getDefaults()
		defaults := &Config{
			App: AppConf{
				UseMemoryStorage: false,
			},
			Logger: LoggerConf{
				LogLevel:        "INFO",
				DisableErrorLog: false,
				DisableWarnLog:  false,
				DisableInfoLog:  false,
				DisableDebugLog: false,
			},
		}

		require.Equal(t, defaults, config)
	})

	t.Run("success read config", func(t *testing.T) {
		config, err := ProvideConfig("testdata/test_config_success.yml")
		expected := &Config{
			App: AppConf{
				Stage:            "local",
				UseMemoryStorage: true,
			},
			HTTP: HTTPConf{
				Host: "localhost",
				Port: "8080",
			},
			Logger: LoggerConf{
				LogLevel:        "INFO",
				DisableErrorLog: false,
				DisableWarnLog:  false,
				DisableInfoLog:  false,
				DisableDebugLog: true,
			},
			Database: DatabaseConf{
				Host:     "localhost",
				Port:     "5432",
				Username: "otus",
				Password: "123456",
				Database: "calendar",
			},
		}

		require.NoError(t, err)
		require.Equal(t, expected, config)
	})

	t.Run("stage required", func(t *testing.T) {
		_, err := ProvideConfig("testdata/test_config_stage_required.yml")
		require.Error(t, err)
	})
}
