package main

import (
	"fmt"
	"log"

	"effective-mobile/internal/config"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("invalid pars config")
	}

	fmt.Printf("%w", cfg)

	// r := chi.NewRouter()
	// r.Route("/api/v1/people", func(r chi.Router) {
	// 	r.Post("/", h.CreatePerson)
	// 	r.Get("/", h.GetPeople)
	// 	r.Put("/{id}", h.UpdatePerson)
	// 	r.Delete("/{id}", h.DeletePerson)
	// })
}
