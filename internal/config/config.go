package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	Env         string
	DatabaseUrl string
}

func MustLoad() Config {
	godotenv.Load()
	
	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT IS REQUIRED")
	}
	env := os.Getenv("ENV")
	if env == "" {
		panic("ENV IS REQUIRED")
	}
	DatabaseUrl := os.Getenv("DATABASE_URL")
	if DatabaseUrl == "" {
		panic("DATABASE_URL IS REQUIRED")
	}

	return Config{
		Port:        port,
		Env:         env,
		DatabaseUrl: DatabaseUrl,
	}
}
