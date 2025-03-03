package config

import (
	"log"
	"time"

	"github.com/joeshaw/envdecode"
	"github.com/joho/godotenv"
)

func readDotEnvFile() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

type Conf struct {
	Server ConfServer
	DB     ConfDB
}

type ConfServer struct {
	Port         int           `env:"SERVER_PORT,required"`
	Address      string        `env:"SERVER_ADDRESS,required"`
	TimeoutRead  time.Duration `env:"SERVER_TIMEOUT_READ,required"`
	TimeoutWrite time.Duration `env:"SERVER_TIMEOUT_WRITE,required"`
	TimeoutIdle  time.Duration `env:"SERVER_TIMEOUT_IDLE,required"`
	Debug        bool          `env:"SERVER_DEBUG,required"`
	JWTSecretKey string        `env:"JWT_SECRET_KEY,required"`
}

type ConfDB struct {
	MongoURI string `env:"MONGO_URI,required"`
	Debug    bool   `env:"DB_DEBUG"`
}

func New() *Conf {
	readDotEnvFile()
	var c Conf
	if err := envdecode.StrictDecode(&c); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}

	if len(c.Server.JWTSecretKey) == 0 {
		log.Fatalf("JWT_SECRET_KEY is required")
	}

	if len(c.Server.JWTSecretKey) < 32 {
		log.Fatalf("JWT_SECRET_KEY must be at least 32 characters")
	}

	return &c
}

func NewDB() *ConfDB {
	readDotEnvFile()
	var c ConfDB
	if err := envdecode.StrictDecode(&c); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}

	return &c
}
