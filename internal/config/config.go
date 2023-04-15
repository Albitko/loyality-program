package config

import (
	"flag"

	"github.com/caarlos0/env/v6"

	"github.com/Albitko/loyalty-program/internal/entities"
)

func New() (entities.Config, error) {
	var cfg entities.Config

	flag.StringVar(&cfg.RunAddress, "a", "localhost:8080", "host and port to listen on")
	flag.StringVar(&cfg.DatabaseURI, "d", "postgresql://localhost:5432/postgres", "database DSN")
	flag.StringVar(&cfg.AccrualSystemAddress, "r", "", "File that stores URL -> ID")
	flag.Parse()

	err := env.Parse(&cfg)

	return cfg, err
}
