package infrastructure

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/mathiasdonoso/hook-replay/internal/domain"
)

type eventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *eventRepository {
	return &eventRepository{db}
}

func (r *eventRepository) Save(event domain.Event) error {
	query := `INSERT INTO events(id, source, path, method, headers, body, received_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`

	id := uuid.New().String()

	_, err := r.db.Exec(
		query,
		id,
		event.Source,
		event.Path,
		event.Method,
		event.Headers,
		event.Body,
		time.Now(),
	)

	return err
}
