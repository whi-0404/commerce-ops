package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTP HTTPConfig `mapstructure:"http"`
}

type HTTPConfig struct {
	Port              int           `mapstructure:"port"`
	ReadTimeout       time.Duration `mapstructure:"read_timeout"`
	ReadHeaderTimeout time.Duration `mapstructure:"read_header_timeout"`
	WriteTimeout      time.Duration `mapstructure:"write_timeout"`
	IdleTimeout       time.Duration `mapstructure:"idle_timeout"`
	ShutdownTimeout   time.Duration `mapstructure:"shutdown_timeout"`
}

func Load() Config {
	return Config{
		HTTP: HTTPConfig{
			Port:              getEnvInt	("HTTP_PORT", 8083),
			ReadTimeout:       getEnvDuration("HTTP_READ_TIMEOUT", 5*time.Second), // read timeout from the client request body
			ReadHeaderTimeout: getEnvDuration("HTTP_READ_HEADER_TIMEOUT", 2*time.Second), // read timeout from the client request header
			WriteTimeout:      getEnvDuration("HTTP_WRITE_TIMEOUT", 10*time.Second), // write timeout for the response to the client
			IdleTimeout:       getEnvDuration("HTTP_IDLE_TIMEOUT", 60*time.Second), // idle timeout for keep-alive connections
			ShutdownTimeout:   getEnvDuration("HTTP_SHUTDOWN_TIMEOUT", 10*time.Second), // timeout for graceful shutdown of the server
		},
	}
}

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	result, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return result
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	result, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}

	return result
}
