package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost             string
	Port                   string
	DBUser                 string
	DBPassword             string
	DBAddress              string
	DBName                 string
	JWTExpirationInSeconds int64
	JWTSecret              string
	BucketName             string
	BucketRegion           string
	S3AccessKey            string
	S3SecretAccessKey      string
}

var Envs = initConfig()

func initConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	return Config{
		PublicHost:             getEnv("PUBLIC_HOST", "localhost"),
		Port:                   getEnv("PORT", "3030"),
		DBUser:                 getEnv("DB_USER", "postgres"),
		DBPassword:             getEnv("DB_PASSWORD", "mypassword"),
		DBAddress:              fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "5432")),
		DBName:                 getEnv("DB_NAME", "postgres"),
		JWTExpirationInSeconds: getEnvAsInt("JWT_EXP", 300),
		JWTSecret:              getEnv("JWT_SECRET", "super-secret-jwt-password"),
		BucketName:             getEnv("BUCKET_NAME", "tireapp-tire-pictures"),
		BucketRegion:           getEnv("BUCKET_REGION", "us-east-1"),
		S3AccessKey:            getEnv("ACCESS_KEY", ""),
		S3SecretAccessKey:      getEnv("SECRET_ACCESS_KEY", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 4)
		if err != nil {
			return fallback
		}
		return i
	}
	return fallback
}
