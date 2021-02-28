package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// LoadConfig - load .env file from given path for local, else will be getting from env var
func LoadConfig() {
	// load .env file from given path for local, else will be getting from env var
	if !strings.EqualFold(os.Getenv("prod"), "true") {
		err := godotenv.Load(".test-env")
		if err != nil {
			panic("Error loading .env file")
		}
	}

	DBConfig = os.Getenv("DB_CONFIG")
	DBConnectionPool, _ = strconv.Atoi(os.Getenv("DB_CONNECTION_POOL"))
	Log, _ = strconv.ParseBool(os.Getenv("LOG"))
	Migrate, _ = strconv.ParseBool(os.Getenv("MIGRATE"))
}
