package helpers

import (
	"github.com/joho/godotenv"
)

func LoadEnv() error {
	return godotenv.Load(".env")

}
