package utils

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var (
	// Port is the PORT environment variable or 8080 if missing.
	// Used to open the tcp listener for our web server.
	Port string
	// MongoDB is the MONGODB environment variable or mongodb://localhost:27017 if missing.
	// Used to connect to the mongodb.
	MongoDB   string
	DbName    string
	JwtSecret string
)

func parse() {
	Port = getDefault("PORT", "8080")
	MongoDB = getDefault("MONGODB", "mongodb://mongodb:27017")
	DbName = getDefault("DB_NAME", "favourites")
	JwtSecret = getDefault("JWT_SECRET_KEY", "")

	log.Printf("• Port=%s\n", Port)
	log.Printf("• MongoDB=%s\n", MongoDB)
	log.Printf("• DbName=%s\n", DbName)
}

// Load loads environment variables that are being used across the whole app.
// Loading from .env file
func Load() {

	envFile := ".env"

	if fileExists(envFile) {
		log.Printf("Loading environment variables from file: %s\n", envFile)

		if err := godotenv.Load(envFile); err != nil {
			panic(fmt.Sprintf("error loading environment variables from [%s]: %v", envFile, err))
		}
	}

	parse()
}

func getDefault(key string, def string) string {
	value := os.Getenv(key)
	if value == "" {
		os.Setenv(key, def)
		value = def
	}

	return value
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
