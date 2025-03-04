package config

import (
	"context"
	"log"

	"github.com/sethvargo/go-envconfig"
)

var (
	Config *config
)

type config struct {
	DEBUG           bool     `env:"DEBUG,default=false"`
	PORT            string   `env:"PORT,default=8000"`
	ALLOWED_ORIGINS []string `env:"ALLOWED_ORIGINS,default=*"`
}

func Load() *config {
	c := config{}
	if err := envconfig.Process(context.TODO(), &c); err != nil {
		log.Fatalf("Failed to load env: %s", err)
	}
	Config = &c
	return &c
}
