package config

import (
	"os"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

type Config struct {
	App      AppConf      `yaml:"app"`
	HTTP     HTTPConf     `yaml:"http"`
	Logger   LoggerConf   `yaml:"logger"`
	Database DatabaseConf `yaml:"db"`
}

type AppConf struct {
	Stage            string `yaml:"stage" validate:"required"`
	UseMemoryStorage bool   `yaml:"useMemoryStorage"`
}

type HTTPConf struct {
	Host string `yaml:"host" validate:"required"`
	Port string `yaml:"port" validate:"required"`
}

type LoggerConf struct {
	LogLevel        string `yaml:"logLevel"`
	DisableErrorLog bool   `yaml:"disableErrorLog"`
	DisableWarnLog  bool   `yaml:"disableWarnLog"`
	DisableInfoLog  bool   `yaml:"disableInfoLog"`
	DisableDebugLog bool   `yaml:"disableDebugLog"`
}

type DatabaseConf struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

func ProvideConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	params := getDefaults()

	err = yaml.Unmarshal(data, params)
	if err != nil {
		return nil, err
	}

	validate := validator.New()
	err = validate.Struct(params)
	if err != nil {
		return nil, err
	}

	return params, nil
}

func getDefaults() *Config {
	return &Config{
		App: AppConf{},
		Logger: LoggerConf{
			LogLevel: "INFO",
		},
	}
}
