package repository

import (
	"context"
	"effective-mobile/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CRUD interface
type PersonRepository interface {
	CreatePerson(ctx context.Context, person *models.PersonEnriched) error
	GetPersons(ctx context.Context, filters models.Filters) ([]models.PersonEnriched, error)
	UpdatePerson(ctx context.Context, id int, update models.PersonUpdate) error
	DeletePerson(ctx context.Context, id int) error
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) CreatePerson(ctx context.Context, person *models.PersonEnriched) error {
	query := `
    INSERT INTO persons 
      (name, surname, patronymic, age, gender, nationality) 
    VALUES 
      (@name, @surname, @patronymic, @age, @gender, @nationality)
    RETURNING id, created_at`

	args := pgx.NamedArgs{
		"name":        person.Name,
		"surname":     person.Surname,
		"patronymic":  person.Patronymic,
		"age":         person.Age,
		"gender":      person.Gender,
		"nationality": person.Nationality,
	}

	return r.pool.QueryRow(ctx, query, args).Scan(
		&person.ID,
		&person.CreatedAt,
	)
}
