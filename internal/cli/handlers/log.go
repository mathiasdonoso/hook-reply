package handlers

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/mathiasdonoso/hook-replay/internal/application"
	"github.com/mathiasdonoso/hook-replay/internal/database"
	"github.com/mathiasdonoso/hook-replay/internal/domain"
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

	printFormatedResponse(events)

	return nil
}

func printFormatedResponse(events []domain.Event) {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintln(w, "Id\t Source\t Path\t Method\t Body\t ReceivedAt")

	bodyMaxLen := 40

	for _, e := range events {
		id := e.Id[:8] + "..."

		body := string(e.Body)

		if len(body) > bodyMaxLen {
			body = body[:bodyMaxLen] + "..."
		}

		fmt.Fprintln(
			w,
			id+"\t",
			e.Source+"\t",
			e.Path+"\t",
			e.Method+"\t",
			body+"\t",
			e.ReceivedAt.Format("2006-01-02 15:04:05")+"\t",
		)
	}

	w.Flush()
}
