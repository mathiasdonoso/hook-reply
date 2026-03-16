package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

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

	// flag: last
	var e domain.Event
	if last {
		e, err = service.Last()
	} else {
		e, err = service.Find(id)
	}
	if err != nil {
		return err
	}

	// flag: target 
	if target != "" {
		if !strings.Contains(target, "://") {
			target = "http://" + target
		}

		target, err := url.Parse(target)
		if err != nil {
			return err
		}

		e.Target = target.String()
	}


	// flags: times & delay
	i := 0
	for i < int(times) {
		if err := makeEventHttpCall(e, delay); err != nil {
			return err
		}
		i = i + 1
	}

	return nil
}

func makeEventHttpCall(e domain.Event, delay uint) error {
	req, err := http.NewRequest(e.Method, e.Target, bytes.NewBuffer(e.Body))
	if err != nil {
		return err
	}

	h := make(map[string][]string)
	err = json.Unmarshal(e.Headers, &h)

	req.Header = h

	time.Sleep(time.Duration(delay) * time.Millisecond)

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	return nil
}
