package config

import (
	"fmt"
	"os"
	"reflect"
)

type Config struct {
	DBDriver   string `env:"DB_DRIVER"`
	DBHost     string `env:"DB_HOST"`
	DBPort     string `env:"DB_PORT"`
	DBName     string `env:"DB_NAME"`
	DBUser     string `env:"DB_USER"`
	DBPassword string `env:"DB_PASSWORD"`
}

func LoadConfig() (*Config, error) {
	config := Config{}

	valueOf := reflect.ValueOf(&config).Elem()
	typeOf := valueOf.Type()

	for i := 0; i < valueOf.NumField(); i++ {
		field := valueOf.Field(i)
		envName := typeOf.Field(i).Tag.Get("env")
		envValue, ok := os.LookupEnv(envName)

		if !ok {
			return nil, fmt.Errorf("%s was not found", envName)
		}

		field.SetString(envValue)
	}

	return &config, nil
}
