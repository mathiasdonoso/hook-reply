package application

import (
	"encoding/json"
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

func (s *eventService) CaptureRequest(r *http.Request, target string) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	headers, err := json.Marshal(r.Header)
	if err != nil {
		return err
	}

	event := domain.Event{
		Source:  r.RemoteAddr,
		Path:    r.RequestURI,
		Method:  r.Method,
		Headers: headers,
		Body:    body,
		Target:  target,
	}

	return s.repo.Save(event)
}

func (s *eventService) ListEvents() ([]domain.Event, error) {
	return s.repo.List()
}

func (s *eventService) Find(id string) (domain.Event, error) {
	return s.repo.Find(id)
}
