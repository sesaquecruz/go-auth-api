package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

type Config struct {
	DBDriver      string `env:"DB_DRIVER"`
	DBHost        string `env:"DB_HOST"`
	DBPort        string `env:"DB_PORT"`
	DBName        string `env:"DB_NAME"`
	DBUser        string `env:"DB_USER"`
	DBPassword    string `env:"DB_PASSWORD"`
	JWTSecret     string `env:"JWT_SECRET"`
	JWTExpSeconds int64  `env:"JWT_EXP_SECONDS"`
}

func LoadConfig() (*Config, error) {
	config := Config{}

	elements := reflect.ValueOf(&config).Elem()
	types := elements.Type()

	for i := 0; i < elements.NumField(); i++ {
		field := elements.Field(i)

		varName := types.Field(i).Tag.Get("env")
		varValue, ok := os.LookupEnv(varName)
		if !ok {
			return nil, fmt.Errorf("%s was not found", varName)
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(varValue)
		case reflect.Int64:
			intValue, err := strconv.Atoi(varValue)
			if err != nil {
				return nil, err
			}
			field.SetInt(int64(intValue))
		default:
			return nil, fmt.Errorf("fail to covert %s", varName)
		}
	}

	return &config, nil
}
