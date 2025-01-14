package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/unidoc/unipdf/v3/common/license"
	"os"
)

func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error memuat file .env")
	}
	return nil
}

func GetEnvVar(key string) string {
	return os.Getenv(key)
}

func SetUnipdfKey(apiKey string) error {
	err := license.SetMeteredKey(apiKey)
	if err != nil {
		return fmt.Errorf("gagal set unipdf key: %v", err)
	}
	return nil
}
