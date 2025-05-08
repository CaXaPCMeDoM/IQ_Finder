package main

import (
	"Name_IQ_Finder/config"
	_ "Name_IQ_Finder/docs"
	"Name_IQ_Finder/internal/app"
	_ "github.com/lib/pq"
	"log"
)

// @title Name IQ Finder API
// @version 1.0
// @description API for Name IQ Finder service
// @host localhost:8080
// @BasePath /api/v1
func main() {
	cfg, err := config.MustLoad()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	app.Run(cfg)
}
