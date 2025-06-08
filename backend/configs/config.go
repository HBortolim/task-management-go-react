package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI          string
	DBName            string
	JWTSecret         string
	Port              string
	TokenExpiryHours  int
	PasswordSaltRound int
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "task_management"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-default-secret-key-change-in-production"
		log.Println("Warning: Using default JWT secret. This should be changed in production.")
	}

	return &Config{
		MongoURI:          mongoURI,
		DBName:            dbName,
		JWTSecret:         jwtSecret,
		Port:              port,
		TokenExpiryHours:  24,
		PasswordSaltRound: 10,
	}
}
