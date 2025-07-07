package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AdamShannag/hookah/internal/auth"
	"github.com/AdamShannag/hookah/internal/condition"
	"github.com/AdamShannag/hookah/internal/config"
	"github.com/AdamShannag/hookah/internal/resolver"
	"github.com/AdamShannag/hookah/internal/server"
	"github.com/AdamShannag/hookah/internal/types"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

func main() {
	templateConfigs, err := parseConfigFile(os.Getenv("CONFIG_PATH"))
	if err != nil {
		log.Fatal(err)
	}

	templates, err := parseTemplates(os.Getenv("TEMPLATES_PATH"))
	if err != nil {
		log.Fatal(err)
	}

	conf := config.New(templateConfigs, templates, auth.NewDefault())

	srv := server.NewServer(conf, condition.NewDefaultEvaluator(resolver.NewPathResolver()))
	done := make(chan bool, 1)
	go gracefulShutdown(srv, done)

	err = srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	<-done
	log.Println("Graceful shutdown complete.")
}

func gracefulShutdown(apiServer *http.Server, done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	done <- true
}

func parseConfigFile(filePath string) ([]types.Template, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var result []types.Template
	if err = json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return result, nil
}

func parseTemplates(dirPath string) (map[string]string, error) {
	dir, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read templates directory: %w", err)
	}

	templates := make(map[string]string)
	for _, entry := range dir {
		name := entry.Name()

		if strings.HasPrefix(name, ".") {
			continue
		}

		fullPath := filepath.Join(dirPath, name)

		info, statErr := os.Stat(fullPath)
		if statErr != nil {
			return nil, fmt.Errorf("failed to stat file: %w", err)
		}
		if info.IsDir() {
			continue
		}

		bytes, readErr := os.ReadFile(fullPath)
		if readErr != nil {
			return nil, fmt.Errorf("failed to read file: %w", readErr)
		}
		templates[name] = string(bytes)
	}

	return templates, nil
}
