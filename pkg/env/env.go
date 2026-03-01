package env

import (
	"os"

	"github.com/joho/godotenv"
)

func Init() {
	godotenv.Load()
}

func Get(key string, fallback string) string {
	env, ok := os.LookupEnv(key)

	if !ok {
		return fallback
	}

	return env
}
