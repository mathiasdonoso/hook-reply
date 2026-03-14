package domain

import "time"

type Event struct {
	Id         string    `json:"id"`
	Source     string    `json:"source"`
	Path       string    `json:"path"`
	Status     string    `json:"status"`
	Headers    string    `json:"headers"`
	Body       []byte    `json:"body"`
	ReceivedAt time.Time `json:"received_at"`
}

type EventRepository interface {
	Save(event Event) error
}
