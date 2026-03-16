package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/mathiasdonoso/hook-replay/internal/application"
	"github.com/mathiasdonoso/hook-replay/internal/database"
	"github.com/mathiasdonoso/hook-replay/internal/infrastructure"
)

func ServeHandler(port uint, forward string) error {
	if forward == "" {
		return fmt.Errorf("invalid URL provided in the --forward flag")
	}

	if !strings.Contains(forward, "://") {
		forward = "http://" + forward
	}

	target, err := url.Parse(forward)
	if err != nil {
		return err
	}

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

	proxy := httputil.NewSingleHostReverseProxy(target)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("incoming request %s %s from %s",
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
		)

		err = service.CaptureRequest(r, forward)
		if err != nil {
			log.Fatal(fmt.Printf("error capturing the request: %s", err.Error()))
		}

		proxy.ServeHTTP(w, r)
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}

	shutdownErr := make(chan error, 1)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdownErr <- server.Shutdown(ctx)
	}()

	log.Printf("Server listening on port %d", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return <-shutdownErr
}
