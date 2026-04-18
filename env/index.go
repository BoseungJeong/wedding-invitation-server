package env

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var AdminPassword string
var AllowOrigin string
var DbPath string

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error: Cannot read .env file")

	}
	AdminPassword = os.Getenv("ADMIN_PASSWORD")
	AllowOrigin = os.Getenv("ALLOW_ORIGIN")
	DbPath = os.Getenv("DB_PATH")
	if DbPath == "" {
		DbPath = "./sql.db"
	}
}
