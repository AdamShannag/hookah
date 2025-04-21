package server

import (
	"encoding/json"
	"io"
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /webhooks/{receiver}", s.WebhookHandler)
	return mux
}

func (s *Server) WebhookHandler(w http.ResponseWriter, r *http.Request) {
	receiver := r.PathValue("receiver")

	for key, _ := range r.URL.Query() {
		r.Header.Set(key, r.URL.Query().Get(key))
	}

	payload, err := io.ReadAll(r.Body)
	if err != nil || len(payload) == 0 {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	templates := s.config.GetConfigTemplates(receiver, r, payload)

	if len(templates) == 0 {
		w.WriteHeader(http.StatusOK)
		return
	}

	var request map[string]interface{}
	if err = json.Unmarshal(payload, &request); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	for _, tmpl := range templates {
		go s.handleTemplate(tmpl, r.Header, request)
	}

	w.WriteHeader(http.StatusOK)
}
