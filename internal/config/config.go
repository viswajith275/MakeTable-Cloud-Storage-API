package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL     string
	Port            string
	SecretKEY       string
	AccessTokenTTL  int
	RefreshTokenTTL int
}

func Load() (*Config, error) {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Failed to load enviornment variables: %v", err)
		return nil, err
	}

	accessTTLstr := os.Getenv("ACCESS_TOKEN_TTL")
	accessTTL, err := strconv.Atoi(accessTTLstr)

	if err != nil {
		log.Fatal("Couldnt convert accesstokenTTl")
		return nil, err
	}
	refreshTTLstr := os.Getenv("REFRESH_TOKEN_TTL")
	refreshTTL, err := strconv.Atoi(refreshTTLstr)

	if err != nil {
		log.Fatal("Couldnt convert refreshtokenTTl")
		return nil, err
	}

	var config *Config = &Config{
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		Port:            os.Getenv("PORT"),
		SecretKEY:       os.Getenv("SECRET_KEY"),
		AccessTokenTTL:  accessTTL,
		RefreshTokenTTL: refreshTTL,
	}

	return config, err

}
