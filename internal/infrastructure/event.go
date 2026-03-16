package infrastructure

import (
	"database/sql"
	// "encoding/json"
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
	q := `INSERT INTO events(id, source, path, method, headers, body, received_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`

	id := uuid.New().String()
	_, err := r.db.Exec(
		q,
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

func (r *eventRepository) List() ([]domain.Event, error) {
	q := `SELECT id, source, path, method, headers, body, received_at FROM events ORDER BY received_at DESC LIMIT 20`
	rows, err := r.db.Query(q)
	if err != nil {
		return []domain.Event{}, err
	}
	defer rows.Close()

	events := []domain.Event{}
	for rows.Next() {
		var e domain.Event
		if err := rows.Scan(
			&e.Id,
			&e.Source,
			&e.Path,
			&e.Method,
			&e.Headers,
			&e.Body,
			&e.ReceivedAt,
		); err != nil {
			return []domain.Event{}, err
		}

		// var body any

		// e.Body = json.Unmarshal(e.Body, body)

		events = append(events, e)
	}

	return events, nil
}
