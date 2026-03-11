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
)

func ServeHandler(port uint, forward string) error {
	if !strings.Contains(forward, "://") {
		forward = "http://" + forward
	}

	target, err := url.Parse(forward)
	if err != nil {
		return fmt.Errorf("invalid forward URL: %w", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: proxy,
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

	log.Println(fmt.Sprintf("Server listening on port %d", port))
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}

	return <-shutdownErr
}
