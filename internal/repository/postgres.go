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
	"go.uber.org/zap"
)

type PostgresRepository struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

func NewPostgresRepository(pool *pgxpool.Pool, logger *zap.Logger) *PostgresRepository {
	return &PostgresRepository{pool: pool, logger: logger}
}

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

func (r *PostgresRepository) CreatePerson(ctx context.Context, person *models.PersonEnriched) error {
	r.logger.Debug("CreatePerson started",
		zap.String("name", person.Name),
		zap.String("surname", person.Surname),
	)

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

	row := r.pool.QueryRow(ctx, query, args)
	if err := row.Scan(&person.ID, &person.CreatedAt); err != nil {
		r.logger.Error("Failed to create person",
			zap.Error(err),
			zap.String("query", query),
			zap.Any("arguments", args),
		)
		return fmt.Errorf("create person: %w", err)
	}

	r.logger.Info("Person created successfully",
		zap.Int("id", person.ID),
		zap.Time("created_at", person.CreatedAt),
	)
	return nil
}

func (r *PostgresRepository) GetPersons(ctx context.Context, filters models.Filters) ([]models.PersonEnriched, error) {
	r.logger.Debug("GetPersons started",
		zap.Any("filters", filters),
		zap.Int("limit", filters.Limit),
		zap.Int("offset", filters.Offset),
	)

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

	r.logger.Debug("Executing SQL query",
		zap.String("query", query.String()),
		zap.Any("args", args),
	)

	rows, err := r.pool.Query(ctx, query.String(), args...)
	if err != nil {
		r.logger.Error("Database query failed",
			zap.Error(err),
			zap.String("query", query.String()),
			zap.Any("arguments", args),
		)
		return nil, fmt.Errorf("database query error: %w", err)
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
			r.logger.Error("Row scanning error",
				zap.Error(err),
				zap.String("query", query.String()),
			)
			return nil, fmt.Errorf("row scan error: %w", err)
		}
		persons = append(persons, p)
	}

	r.logger.Info("Successfully fetched persons",
		zap.Int("count", len(persons)),
		zap.Int("limit", filters.Limit),
		zap.Int("page", filters.Offset/filters.Limit+1),
	)

	return persons, nil
}

func (r *PostgresRepository) UpdatePerson(ctx context.Context, id int, upd models.PersonUpdate) error {
	r.logger.Debug("UpdatePerson started",
		zap.Any("PersonUpdate", upd),
	)

	var (
		setClauses []string
		args       []interface{}
		argPos     = 1
	)

	if upd.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argPos))
		args = append(args, *upd.Name)
		argPos++
	}
	if upd.Surname != nil {
		setClauses = append(setClauses, fmt.Sprintf("surname = $%d", argPos))
		args = append(args, *upd.Surname)
		argPos++
	}
	if upd.Patronymic != nil {
		setClauses = append(setClauses, fmt.Sprintf("patronymic = $%d", argPos))
		args = append(args, *upd.Patronymic)
		argPos++
	}
	if upd.Age != nil {
		setClauses = append(setClauses, fmt.Sprintf("age = $%d", argPos))
		args = append(args, *upd.Age)
		argPos++
	}
	if upd.Gender != nil {
		setClauses = append(setClauses, fmt.Sprintf("gender = $%d", argPos))
		args = append(args, *upd.Gender)
		argPos++
	}
	if upd.Nationality != nil {
		setClauses = append(setClauses, fmt.Sprintf("nationality = $%d", argPos))
		args = append(args, *upd.Nationality)
		argPos++
	}

	if len(setClauses) == 0 {
		return errors.New("no fields to update")
	}

	query := fmt.Sprintf(
		"UPDATE persons SET %s WHERE id = $%d",
		strings.Join(setClauses, ", "),
		argPos,
	)
	args = append(args, id)

	r.logger.Debug("Executing SQL query",
		zap.String("query", query),
		zap.Any("args", args),
	)

	_, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		r.logger.Error("Database query failed",
			zap.Error(err),
			zap.String("query", query),
			zap.Any("arguments", args),
		)
		return fmt.Errorf("database query error: %w", err)
	}

	return nil
}
func (r *PostgresRepository) DeletePerson(ctx context.Context, id int) error {
	r.logger.Debug("DeletePerson started",
		zap.Int("id", id),
	)

	query := "DELETE FROM persons WHERE id = $1"

	r.logger.Debug("Executing delete query",
		zap.String("query", query),
		zap.Int("id", id),
	)

	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("Failed to delete person",
			zap.Error(err),
			zap.String("query", query),
			zap.Int("id", id),
		)
		return fmt.Errorf("error deleting person: %w", err)
	}

	rowsAffected := tag.RowsAffected()
	if rowsAffected == 0 {
		r.logger.Warn("Person not found",
			zap.Int("id", id),
			zap.String("query", query),
		)
		return fmt.Errorf("person with id %d not found", id)
	}

	r.logger.Info("Successfully deleted person",
		zap.Int("id", id),
		zap.Int64("rows_affected", rowsAffected),
	)

	return nil
}
