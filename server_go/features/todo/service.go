package todo

import (
	"context"
	"database/sql"
	"strings"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) list(ctx context.Context) ([]Todo, error) {
	rows, err := s.db.QueryContext(
		ctx,
		"SELECT * FROM todo ORDER BY id ASC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := []Todo{}
	for rows.Next() {
		var t Todo
		if err := rows.Scan(
			&t.ID, &t.Task, &t.Status, &t.Priority, &t.UpdatedAt, &t.CreatedAt,
		); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}

	return todos, rows.Err()
}

func (s *Service) create(ctx context.Context, task string, priority string) (Todo, error) {
	var t Todo

	if priority == "" {
		priority = "medium"
	}

	err := s.db.QueryRowContext(
		ctx,
		"INSERT INTO todo (task, priority, status, updated_at) VALUES (?, ?, 'todo', CURRENT_TIMESTAMP) "+
			"RETURNING id, task, status, priority, updated_at, created_at",
		task,
		priority,
	).Scan(
		&t.ID, &t.Task, &t.Status, &t.Priority, &t.UpdatedAt, &t.CreatedAt,
	)
	if err != nil {
		return Todo{}, err
	}

	return t, nil
}

func (s *Service) update(ctx context.Context, id int64, task *string, priority *string) (Todo, error) {
	var t Todo

	err := s.db.QueryRowContext(
		ctx,
		"UPDATE todo SET task = COALESCE(?, task), priority = COALESCE(?, priority), updated_at = CURRENT_TIMESTAMP WHERE id = ? "+
			"RETURNING id, task, status, priority, updated_at, created_at",
		task,
		priority,
		id,
	).Scan(
		&t.ID, &t.Task, &t.Status, &t.Priority, &t.UpdatedAt, &t.CreatedAt,
	)
	if err != nil {
		return Todo{}, err
	}

	return t, nil
}

func (s *Service) updateStatus(ctx context.Context, id int64, status string) (Todo, error) {
	var t Todo

	err := s.db.QueryRowContext(
		ctx,
		"UPDATE todo SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? "+
			"RETURNING id, task, status, priority, updated_at, created_at",
		status,
		id,
	).Scan(
		&t.ID, &t.Task, &t.Status, &t.Priority, &t.UpdatedAt, &t.CreatedAt,
	)
	if err != nil {
		return Todo{}, err
	}

	return t, nil
}

func (s *Service) delete(ctx context.Context, ids []int64) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	// Create a string of placeholders for the IN clause
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	query := "DELETE FROM todo WHERE id IN (" + strings.Join(placeholders, ", ") + ")"
	res, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}
