package server

import (
	"fmt"
	"github.com/AdamShannag/hookah/internal/types"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Server struct {
	port   int
	config *types.Config
}

func NewServer(config *types.Config) *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	newServer := &Server{
		port:   port,
		config: config,
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
