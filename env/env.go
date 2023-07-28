package env

import (
	"os"

	"github.com/joho/godotenv"
)

var err = godotenv.Load(".env")
var envs = map[string]string{}

func Getenv(key string) string {
	if envs[key] != "" {
		return envs[key]
	}

	envs[key] = os.Getenv(key)

	return envs[key]
}
