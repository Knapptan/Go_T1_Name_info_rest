package main

import (
	"context"
	"net/http"

	_ "effective-mobile/docs"
	"effective-mobile/internal/clients"
	"effective-mobile/internal/config"
	"effective-mobile/internal/handler"
	"effective-mobile/internal/repository"
	"effective-mobile/internal/service"

	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

// @title Effective Mobile API
// @version 1.0
// @description REST API для работы с данными людей

// @contact.name API Support
// @contact.email support@effective-mobile.ru

// @host localhost:8080
// @BasePath /api/v1
func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Config load error", zap.Error(err))
	}

	pool, err := repository.NewPostgresDB(context.Background(), cfg.DatabaseURL())
	if err != nil {
		logger.Fatal("DB connection error", zap.Error(err))
	}
	defer pool.Close()

	repo := repository.NewPostgresRepository(pool, logger)
	enricher := clients.NewEnrichmentClient(logger)
	svc := service.NewPersonService(repo, enricher, logger)
	h := handler.NewHandler(svc, logger)

	r := chi.NewRouter()
	r.Get("/swagger/*", httpSwagger.WrapHandler)
	r.Route("/api/v1/people", func(r chi.Router) {
		r.Post("/", h.CreatePerson)
		r.Get("/", h.GetPeople)
		r.Patch("/{id}", h.UpdatePerson)
		r.Delete("/{id}", h.DeletePerson)
	})

	logger.Info("Starting server",
		zap.String("port", cfg.Port),
		zap.String("environment", cfg.Environment),
	)

	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		logger.Fatal("Server failed", zap.Error(err))
	}
}
