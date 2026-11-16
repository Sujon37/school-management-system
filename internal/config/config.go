package config

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AdminConfig struct {
	Email    string `env:"ADMIN_EMAIL" env-required:"true"`
	Password string `env:"ADMIN_PASSWORD" env-required:"true"`
}

type HTTPServer struct {
	Addr string
}

type Config struct {
	Env         string
	StoragePath string
	HTTPServer  HTTPServer
	Admin       AdminConfig
	JWTSecret   string
}

func MustLoad() *Config {
	var configPath string

	flag.StringVar(&configPath, "config", "", "path to configuration file")
	flag.Parse()

	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}

	if configPath == "" {
		configPath = ".env"
	}

	if err := godotenv.Load(configPath); err != nil {
		log.Println("Warning: .env file not found, reading directly from system environment variables.")
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	storagePath := os.Getenv("STORAGE_PATH")
	if storagePath == "" {
		storagePath = "storage/school.db"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	admin := AdminConfig{
		Email:    os.Getenv("ADMIN_EMAIL"),
		Password: os.Getenv("ADMIN_PASSWORD"),
	}

	if admin.Email == "" || admin.Password == "" {
		log.Fatal("ADMIN_EMAIL and ADMIN_PASSWORD are required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")

	if len(jwtSecret) < 32 {
		log.Fatal("JWT_SECRET must be at least 32 bytes")
	}

	return &Config{
		Env:         env,
		StoragePath: storagePath,
		HTTPServer: HTTPServer{
			Addr: ":" + port,
		},
		Admin: admin,
		JWTSecret: jwtSecret,
	}
}
