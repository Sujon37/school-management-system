
package config

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type HTTPServer struct {
	Addr string
}

type Config struct {
	Env         string
	StoragePath string
	HTTPServer  HTTPServer
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

	return &Config{
		Env:         env,
		StoragePath: storagePath,
		HTTPServer: HTTPServer{
			Addr: ":" + port,
		},
	}
}
