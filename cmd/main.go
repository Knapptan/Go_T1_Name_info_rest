package main

import (
	"context"
	"log"

	"effective-mobile/internal/clients"
	"effective-mobile/internal/config"
	"effective-mobile/internal/models"
	"effective-mobile/internal/repository"
	"effective-mobile/internal/service"

	"github.com/go-chi/chi"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("invalid pars config")
	}

	ctx := context.Background()
	pool, err := repository.NewPostgresDB(ctx, cfg.DatabaseURL())
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}
	defer pool.Close()

	repo := repository.NewPostgresRepository(pool)

	enricher := clients.NewEnrichmentClient()

	personService := service.NewPersonService(repo, enricher)

	patronymic := "Vasilevich"

	personInput := models.PersonInput{
		Name:       "Dmitriy",
		Surname:    "Ushakov",
		Patronymic: &patronymic,
	}

	person, err := personService.CreatePersonEnriched(ctx, personInput)
	if err != nil {
		log.Fatalf("CreatePersonEnriched failed: %v", err)
	}

	if err := repo.CreatePerson(ctx, person); err != nil {
		log.Fatalf("CreatePerson failed: %v", err)
	}

	r := chi.NewRouter()
	r.Route("/api/v1/people", func(r chi.Router) {
		// 	r.Post("/", h.CreatePerson)
		// 	r.Get("/", h.GetPeople)
		// 	r.Put("/{id}", h.UpdatePerson)
		// 	r.Delete("/{id}", h.DeletePerson)
	})
}
