package application

import (
	"fmt"
	"io"
	"net/http"

	"github.com/mathiasdonoso/hook-replay/internal/domain"
)

type eventService struct {
	repo domain.EventRepository
}

func NewEventService(repo domain.EventRepository) *eventService {
	return &eventService{repo: repo}
}

func (s *eventService) Capture(r *http.Request) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("err reading r.Body: %+v\n", err)
		return err
	}

	event := domain.Event{
		Source:  r.RemoteAddr,
		Path:    r.RequestURI,
		Status:  "200",
		Headers: "",
		Body:    body,
	}

	return s.repo.Save(event)
}
