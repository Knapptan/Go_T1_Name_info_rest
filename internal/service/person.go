package service

import (
	"context"
	"effective-mobile/internal/clients"
	"effective-mobile/internal/models"
	"effective-mobile/internal/repository"
	"fmt"
	"log"
)

type PersonService struct {
	repo     repository.PersonRepository
	enricher *clients.EnrichmentClient
}

func NewPersonService(repo repository.PersonRepository, enricher *clients.EnrichmentClient) *PersonService {
	return &PersonService{
		repo:     repo,
		enricher: enricher,
	}
}

func (s *PersonService) CreatePersonEnriched(ctx context.Context, input models.PersonInput) (*models.PersonEnriched, error) {
	log.Printf("Starting enrichment for: %s %s", input.Name, input.Surname)
	// Обогащение данных через внешние API
	age, err := s.enricher.GetAge(input.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to enricher.GetAge: %w", err)
	}

	gender, err := s.enricher.GetGender(input.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to enricher.GetGender: %w", err)
	}

	nationality, err := s.enricher.GetNationality(input.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to enricher.GetNationality: %w", err)
	}

	log.Printf("Enrichment results - Age: %d, Gender: %s, Nationality: %s",
		age, gender, nationality)

	person := models.PersonEnriched{
		Name:        input.Name,
		Surname:     input.Surname,
		Patronymic:  input.Patronymic,
		Age:         age,
		Gender:      gender,
		Nationality: nationality,
	}

	if err := s.repo.CreatePerson(ctx, &person); err != nil {
		return nil, fmt.Errorf("failed repo.CreatePerson: %w", err)
	}

	log.Printf("Successfully created person with ID: %d", person.ID)

	return &person, nil
}
