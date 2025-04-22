package server

import (
	"fmt"
	"github.com/AdamShannag/hookah/internal/condition"
	"github.com/AdamShannag/hookah/internal/config"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Server struct {
	port      int
	config    *config.Config
	evaluator condition.Evaluator
}

func NewServer(config *config.Config, evaluator condition.Evaluator) *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	newServer := &Server{
		port:      port,
		config:    config,
		evaluator: evaluator,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", newServer.port),
		Handler:      newServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
