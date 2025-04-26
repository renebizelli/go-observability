package configs

import (
	"renebizelli/go/observability/SystemB/utils"

	"github.com/spf13/viper"
)

type Config struct {
	WEBSERVER_PORT              int
	SERVICES_TIMEOUT            int
	VIACEP_URL                  string
	WEATHERAPI_URL              string
	WEATHERAPI_KEY              string
	OTEL_SERVICE_NAME           string
	OTEL_EXPORTER_OTLP_ENDPOINT string
}

func LoadConfig(path string) *Config {

	var cfg *Config
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	utils.PanicIfError(err, "Load config file error")

	err = viper.Unmarshal(&cfg)
	utils.PanicIfError(err, "Unmarshal error")

	return cfg
}
