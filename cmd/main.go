package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"effective-mobile/internal/config"
	"effective-mobile/internal/models"
	"effective-mobile/internal/repository"

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

	Patronymic := "Auto"

	person := &models.PersonEnriched{
		Name:        "Test",
		Surname:     "User",
		Patronymic:  &Patronymic,
		Age:         30,
		Gender:      "male",
		Nationality: "US",
	}

	if err := repo.CreatePerson(ctx, person); err != nil {
		log.Fatalf("CreatePerson failed: %v", err)
	}

	fmt.Printf("Created person: ID=%d, CreatedAt=%s\n", person.ID, person.CreatedAt.Format(time.RFC3339))

	r := chi.NewRouter()
	r.Route("/api/v1/people", func(r chi.Router) {
		// 	r.Post("/", h.CreatePerson)
		// 	r.Get("/", h.GetPeople)
		// 	r.Put("/{id}", h.UpdatePerson)
		// 	r.Delete("/{id}", h.DeletePerson)
	})
}
