package repo

import (
	"database/sql"
	"fmt"
	"time"

	"Name_IQ_Finder/internal/entity"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) Create(person *entity.Person) (int64, error) {
	query := `
		INSERT INTO persons (name, surname, patronymic, age, gender, nationality, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	now := time.Now().Format(time.RFC3339)
	person.CreatedAt = now
	person.UpdatedAt = now

	var id int64
	err := r.db.QueryRow(
		query,
		person.Name,
		person.Surname,
		person.Patronymic,
		person.Age,
		person.Gender,
		person.Nationality,
		person.CreatedAt,
		person.UpdatedAt,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create person: %w", err)
	}

	return id, nil
}

func (r *PostgresRepository) GetByID(id int64) (*entity.Person, error) {
	query := `
		SELECT id, name, surname, patronymic, age, gender, nationality, created_at, updated_at
		FROM persons
		WHERE id = $1
	`

	var person entity.Person
	err := r.db.QueryRow(query, id).Scan(
		&person.ID,
		&person.Name,
		&person.Surname,
		&person.Patronymic,
		&person.Age,
		&person.Gender,
		&person.Nationality,
		&person.CreatedAt,
		&person.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("person not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get person: %w", err)
	}

	return &person, nil
}

func (r *PostgresRepository) GetAll(filter map[string]interface{}, page, limit int) ([]*entity.Person, int, error) {
	query := `
		SELECT id, name, surname, patronymic, age, gender, nationality, created_at, updated_at
		FROM persons
	`

	countQuery := `
		SELECT COUNT(*)
		FROM persons
	`

	whereClause := ""
	args := []interface{}{}
	argIndex := 1

	if len(filter) > 0 {
		whereClause = " WHERE "
		for key, value := range filter {
			if argIndex > 1 {
				whereClause += " AND "
			}
			whereClause += fmt.Sprintf("%s = $%d", key, argIndex)
			args = append(args, value)
			argIndex++
		}
	}

	offset := (page - 1) * limit
	paginationClause := fmt.Sprintf(" ORDER BY id LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	var totalCount int
	err := r.db.QueryRow(countQuery+whereClause, args[:len(args)-2]...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count persons: %w", err)
	}

	rows, err := r.db.Query(query+whereClause+paginationClause, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get persons: %w", err)
	}
	defer rows.Close()

	var persons []*entity.Person
	for rows.Next() {
		var person entity.Person
		err := rows.Scan(
			&person.ID,
			&person.Name,
			&person.Surname,
			&person.Patronymic,
			&person.Age,
			&person.Gender,
			&person.Nationality,
			&person.CreatedAt,
			&person.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan person: %w", err)
		}
		persons = append(persons, &person)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating persons: %w", err)
	}

	return persons, totalCount, nil
}

func (r *PostgresRepository) Update(person *entity.Person) error {
	query := `
		UPDATE persons
		SET name = $1, surname = $2, patronymic = $3, age = $4, gender = $5, nationality = $6, updated_at = $7
		WHERE id = $8
	`

	person.UpdatedAt = time.Now().Format(time.RFC3339)

	result, err := r.db.Exec(
		query,
		person.Name,
		person.Surname,
		person.Patronymic,
		person.Age,
		person.Gender,
		person.Nationality,
		person.UpdatedAt,
		person.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update person: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("person not found")
	}

	return nil
}

func (r *PostgresRepository) Delete(id int64) error {
	query := `
		DELETE FROM persons
		WHERE id = $1
	`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete person: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("person not found")
	}

	return nil
}
