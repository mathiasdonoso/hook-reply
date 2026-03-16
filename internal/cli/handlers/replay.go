package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/mathiasdonoso/hook-replay/internal/application"
	"github.com/mathiasdonoso/hook-replay/internal/database"
	"github.com/mathiasdonoso/hook-replay/internal/domain"
	"github.com/mathiasdonoso/hook-replay/internal/infrastructure"
)

// todo: use a different struct to represent handler params & flags?
func ReplayHandler(id string, last bool, times uint, delay uint, target string) error {
	connConfig, err := database.NewConfig()
	if err != nil {
		return err
	}

	conn, err := database.NewConnection(connConfig)
	if err != nil {
		return err
	}

	defer conn.Close()

	eventRepo := infrastructure.NewEventRepository(conn.DB())
	service := application.NewEventService(eventRepo)

	e, err := service.Find(id)
	if err != nil {
		return err
	}

	return makeEventHttpCall(e)
}

func makeEventHttpCall(e domain.Event) error {
	req, err := http.NewRequest(e.Method, e.Target, bytes.NewBuffer(e.Body))
	if err != nil {
		return err
	}

	h := make(map[string][]string)
	err = json.Unmarshal(e.Headers, &h)

	req.Header = h

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	return nil
}
