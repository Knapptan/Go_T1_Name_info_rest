package service

import (
	"context"
	"effective-mobile/internal/clients"
	"effective-mobile/internal/models"
	"effective-mobile/internal/repository"
	"fmt"
	"log"

	"go.uber.org/zap"
)

type PersonService struct {
	repo     repository.PersonRepository
	enricher *clients.EnrichmentClient
	logger   *zap.Logger
}

func NewPersonService(repo repository.PersonRepository, enricher *clients.EnrichmentClient, logger *zap.Logger) *PersonService {
	return &PersonService{
		repo:     repo,
		enricher: enricher,
		logger:   logger.With(zap.String("component", "person_service")),
	}
}

func (s *PersonService) CreatePersonEnriched(ctx context.Context, input models.PersonInput) (*models.PersonEnriched, error) {
	log.Printf("Starting enrichment for: %s %s", input.Name, input.Surname)

	s.logger.Debug("Starting person enrichment",
		zap.String("name", input.Name),
		zap.String("surname", input.Surname),
	)

	age, err := s.enricher.GetAge(input.Name)
	if err != nil {
		s.logger.Error("Failed to get age",
			zap.String("name", input.Name),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get age: %w", err)
	}

	gender, err := s.enricher.GetGender(input.Name)
	if err != nil {
		s.logger.Error("Failed to get gender",
			zap.String("name", input.Name),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get gender: %w", err)
	}

	nationality, err := s.enricher.GetNationality(input.Name)
	if err != nil {
		s.logger.Error("Failed to get nationality",
			zap.String("name", input.Name),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get nationality: %w", err)
	}

	s.logger.Info("Enrichment completed",
		zap.Int("age", age),
		zap.String("gender", gender),
		zap.String("nationality", nationality),
	)

	person := models.PersonEnriched{
		Name:        input.Name,
		Surname:     input.Surname,
		Patronymic:  input.Patronymic,
		Age:         age,
		Gender:      gender,
		Nationality: nationality,
	}

	if err := s.repo.CreatePerson(ctx, &person); err != nil {
		s.logger.Error("Failed to create person",
			zap.Any("person", person),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to create person: %w", err)
	}

	s.logger.Info("Person created successfully",
		zap.Int("id", person.ID),
		zap.Time("created_at", person.CreatedAt),
	)

	return &person, nil
}

func (s *PersonService) GetPersons(ctx context.Context, filters models.Filters) ([]models.PersonEnriched, error) {
	s.logger.Debug("Getting persons with filters",
		zap.Any("filters", filters),
	)

	persons, err := s.repo.GetPersons(ctx, filters)
	if err != nil {
		s.logger.Error("Failed to get persons",
			zap.Any("filters", filters),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get persons: %w", err)
	}

	s.logger.Info("Persons fetched",
		zap.Int("count", len(persons)),
	)

	return persons, nil
}

func (s *PersonService) DeletePerson(ctx context.Context, id int) error {
	s.logger.Debug("Deleting person",
		zap.Int("id", id),
	)

	if err := s.repo.DeletePerson(ctx, id); err != nil {
		s.logger.Error("Failed to delete person",
			zap.Int("id", id),
			zap.Error(err),
		)
		return fmt.Errorf("failed to delete person: %w", err)
	}

	s.logger.Info("Person deleted",
		zap.Int("id", id),
	)

	return nil
}

func (s *PersonService) UpdatePerson(ctx context.Context, id int, update models.PersonUpdate) error {
	s.logger.Debug("Updating person",
		zap.Int("id", id),
		zap.Any("update", update),
	)

	if err := s.repo.UpdatePerson(ctx, id, update); err != nil {
		s.logger.Error("Failed to update person",
			zap.Int("id", id),
			zap.Any("update", update),
			zap.Error(err),
		)
		return fmt.Errorf("failed to update person: %w", err)
	}

	s.logger.Info("Person updated",
		zap.Int("id", id),
	)

	return nil
}
