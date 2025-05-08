package app

import (
	"Name_IQ_Finder/config"
	"Name_IQ_Finder/internal/controller/http"
	"Name_IQ_Finder/internal/infrastructure/api"
	"Name_IQ_Finder/internal/infrastructure/repo"
	"Name_IQ_Finder/internal/logger"
	"Name_IQ_Finder/internal/usecase"
	"database/sql"
	"log"
	"os"
)

func Run(cfg *config.Config) {
	appLogger := logger.New(cfg.Logger.Level)
	appLogger.Info("Starting Name IQ Finder service")

	db, err := sql.Open("postgres", cfg.Database.GetDSN())
	if err != nil {
		appLogger.Fatal("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		appLogger.Fatal("Failed to ping database: %v", err)
	}
	appLogger.Info("Connected to database")

	personRepo := repo.NewPostgresRepository(db)

	externalClient := api.NewExternalClient()

	personUseCase := usecase.NewPersonUseCase(personRepo, externalClient, log.New(os.Stdout, "", log.LstdFlags))

	router := http.NewRouter(personUseCase)

	appLogger.Info("Starting server on port %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		appLogger.Fatal("Failed to start server: %v", err)
	}
}
