package main

import (
	"flag"
	"log"

	_ "github.com/lib/pq"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/config"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/app"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./config/config.yml", "Path to configuration file")
}

func main() {
	flag.Parse()

	// Configuration
	cfg, err := config.NewConfig(configFile)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
