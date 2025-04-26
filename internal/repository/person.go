package repository

import (
	"context"
	"effective-mobile/internal/models"
)

// CRUD interface
type PersonRepository interface {
	CreatePerson(ctx context.Context, person *models.PersonEnriched) error
	GetPersons(ctx context.Context, filters models.Filters) ([]models.PersonEnriched, error)
	UpdatePerson(ctx context.Context, id int, update models.PersonUpdate) error
	DeletePerson(ctx context.Context, id int) error
}
