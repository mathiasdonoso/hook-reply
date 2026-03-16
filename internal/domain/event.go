package domain

import "time"

type Event struct {
	Id         string    `json:"id"`
	Source     string    `json:"source"`
	Path       string    `json:"path"`
	Method     string    `json:"method"`
	Headers    []byte    `json:"headers"`
	Body       []byte    `json:"body"`
	Target     string    `json:"target"`
	ReceivedAt time.Time `json:"received_at"`
}

type EventRepository interface {
	Save(event Event) error
	List() ([]Event, error)
	Find(id string) (Event, error)
	Last() (Event, error)
}
