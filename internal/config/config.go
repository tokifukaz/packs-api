package config

import (
	"fmt"
	"os"
	"strings"

	"packs-api/internal/store"
)

type Config struct {
	Addr                   string
	PathPrefix             string
	MongoDB                store.NoSQLStore
	AllowedOrigins         []string
	SkipHealthCheckLogging bool
}

func NewConfig(addr string) (*Config, error) {
	cfg := new(Config)
	cfg.Addr = addr
	cfg.PathPrefix = "/api"

	ao, ok := os.LookupEnv("ALLOWED_ORIGINS")
	if !ok {
		cfg.AllowedOrigins = []string{"*"}
	} else {
		cfg.AllowedOrigins = strings.Split(ao, ",")
	}

	m, err := getMongoDB()
	if err != nil {
		return nil, err
	}

	cfg.MongoDB = m

	return cfg, nil
}

func getMongoDB() (store.NoSQLStore, error) {
	mongoURI := os.Getenv("MONGODB_URI")
	mongoDBName := os.Getenv("MONGODB_DATABASE_NAME")
	mongoCertPath := os.Getenv("MONGODB_CERT_PATH")

	mongoDB, err := store.NewMongoDB(mongoURI, mongoDBName, mongoCertPath)
	if err != nil {
		return nil, fmt.Errorf("error connecting to MongoDB: %s", err)
	}

	return mongoDB, nil
}
