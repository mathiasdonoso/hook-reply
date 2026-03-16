package handlers

import (
	"fmt"

	"github.com/mathiasdonoso/hook-replay/internal/application"
	"github.com/mathiasdonoso/hook-replay/internal/database"
	"github.com/mathiasdonoso/hook-replay/internal/infrastructure"
)

func LogHandler() error {
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
	events, err := service.ListEvents()
	if err != nil {
		return err
	}

	for _, e := range events {
		fmt.Printf("%+v\n", e.Body)
	}

	return nil
}
