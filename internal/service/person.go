package service

import (
	"context"
	"effective-mobile/internal/models"
	"effective-mobile/internal/repository"
)

type PersonService struct {
	repo repository.PersonRepository
}

func NewPersonService(repo repository.PersonRepository) *PersonService {
	return &PersonService{repo: repo}
}

func (s *PersonService) CreatePersonEnriched(ctx context.Context, input models.PersonInput) (*models.PersonEnriched, error) {
	// Обогащение данных через внешние API
	enriched := models.Enrichment{}
	// TODO заполеннеие полей через API

	person := models.PersonEnriched{
		Name:        input.Name,
		Surname:     input.Surname,
		Patronymic:  &input.Patronymic,
		Age:         enriched.Age,
		Gender:      enriched.Gender,
		Nationality: enriched.Nationality,
	}

	if err := s.repo.CreatePerson(ctx, &person); err != nil {
		return nil, err
	}

	return &person, nil
}
