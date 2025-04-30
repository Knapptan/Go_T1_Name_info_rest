package repository

import (
	"context"
	"effective-mobile/internal/models"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresDB(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse conn string: %v", err)
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %v", err)
	}

	return pool, nil
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

func (r *PostgresRepository) GetPersons(ctx context.Context, filters models.Filters) ([]models.PersonEnriched, error) {
	query := strings.Builder{}
	query.WriteString(`
	    SELECT
      id, name, surname, patronymic,
      age, gender, nationality, created_at
    FROM persons
    WHERE 1=1`)

	args := make([]interface{}, 0)
	argsCounter := 1

	if filters.Name != "" {
		query.WriteString(fmt.Sprintf(" AND name = $%d", argsCounter))
		args = append(args, filters.Name)
		argsCounter++
	}

	if filters.Gender != "" {
		query.WriteString(fmt.Sprintf(" AND gender = $%d", argsCounter))
		args = append(args, filters.Gender)
		argsCounter++
	}

	query.WriteString(fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argsCounter, argsCounter+1))
	args = append(args, filters.Limit, filters.Offset)

	rows, err := r.pool.Query(ctx, query.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var persons []models.PersonEnriched

	for rows.Next() {
		var p models.PersonEnriched
		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Surname,
			&p.Patronymic,
			&p.Age,
			&p.Gender,
			&p.Nationality,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		persons = append(persons, p)
	}

	return persons, nil
}

func (r *PostgresRepository) UpdatePerson(ctx context.Context, id int, update models.PersonUpdate) error {
	query := strings.Builder{}
	query.WriteString("UPDATE persons SET")

	args := make([]interface{}, 0)
	argCounter := 1

	if update.Age != nil {
		query.WriteString(fmt.Sprintf(" age = $%d", argCounter))
		args = append(args, *update.Age)
		argCounter++
	}

	if update.Gender != nil {
		query.WriteString(fmt.Sprintf(" gender = $%d", argCounter))
		args = append(args, *update.Gender)
		argCounter++
	}

	if update.Nationality != nil {
		query.WriteString(fmt.Sprintf(" nationality = $%d", argCounter))
		args = append(args, *update.Nationality)
		argCounter++
	}

	if len(args) == 0 {
		return errors.New("no fields to update")
	}

	queryStr := strings.TrimSuffix(query.String(), ",") + fmt.Sprintf(" WHERE id = $%d", argCounter)

	args = append(args, id)

	_, err := r.pool.Exec(ctx, queryStr, args...)
	return err
}

func (r *PostgresRepository) DeletePerson(ctx context.Context, id int) error {
	tag, err := r.pool.Exec(ctx,
		"DELETE FROM persons WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("error deleting person: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("person with id %d not found", id)
	}
	return nil
}
