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
	q := `INSERT INTO events(id, source, path, method, headers, body, target, received_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`

	id := uuid.New().String()
	_, err := r.db.Exec(
		q,
		id,
		event.Source,
		event.Path,
		event.Method,
		event.Headers,
		event.Body,
		event.Target,
		time.Now(),
	)

	return err
}

func (r *eventRepository) List() ([]domain.Event, error) {
	q := `SELECT id, source, path, method, headers, body, target, received_at FROM events ORDER BY received_at DESC LIMIT 20;`
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
			&e.Target,
			&e.ReceivedAt,
		); err != nil {
			return []domain.Event{}, err
		}

		events = append(events, e)
	}

	return events, nil
}

func (r *eventRepository) Find(id string) (domain.Event, error) {
	var event domain.Event
	q := `SELECT id, source, path, method, headers, body, target, received_at FROM events WHERE id LIKE $1;`

	if err := r.db.QueryRow(q, id+"%").Scan(
		&event.Id,
		&event.Source,
		&event.Path,
		&event.Method,
		&event.Headers,
		&event.Body,
		&event.Target,
		&event.ReceivedAt,
	); err != nil {
		return event, err
	}

	return event, nil
}
