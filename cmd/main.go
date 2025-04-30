package main

import (
	"context"
	"log"
	"net/http"

	"effective-mobile/internal/clients"
	"effective-mobile/internal/config"
	"effective-mobile/internal/handler"
	"effective-mobile/internal/repository"
	"effective-mobile/internal/service"

	"github.com/go-chi/chi"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load error: %v", err)
	}

	pool, err := repository.NewPostgresDB(context.Background(), cfg.DatabaseURL())
	if err != nil {
		log.Fatalf("db connection error: %v", err)
	}
	defer pool.Close()

	repo := repository.NewPostgresRepository(pool)

	enricher := clients.NewEnrichmentClient()

	svc := service.NewPersonService(repo, enricher)
	h := handler.NewHandler(svc)

	r := chi.NewRouter()
	r.Route("/api/v1/people", func(r chi.Router) {
		r.Post("/", h.CreatePerson)
		// 	r.Get("/", h.GetPeople)
		// 	r.Put("/{id}", h.UpdatePerson)
		// 	r.Delete("/{id}", h.DeletePerson)
	})

	log.Printf("starting server on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
