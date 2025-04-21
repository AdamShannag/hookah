package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/AdamShannag/hookah/internal/condition"
	"github.com/AdamShannag/hookah/internal/render"
	"github.com/AdamShannag/hookah/internal/types"
	"log"
	"net/http"
)

func (s *Server) handleTemplate(tmpl types.Template, headers http.Header, body map[string]any) {
	eventType, err := extractEventType(tmpl, headers, body)
	if err != nil {
		log.Printf("[Template] %v", err)
		return
	}

	events := tmpl.Events.GetEvents(eventType)
	if len(events) == 0 {
		log.Printf("[Template] No matching events found for type: %v", eventType)
		return
	}

	for _, evt := range events {
		go s.processEvent(evt, headers, body)
	}
}

func (s *Server) processEvent(evt types.Event, headers http.Header, body map[string]any) {
	ok, err := condition.Evaluate(evt.Conditions, headers, body)
	if err != nil {
		log.Printf("[Condition] Evaluation error: %v", err)
		return
	}
	if !ok {
		log.Println("[Condition] Not met, skipping event")
		return
	}

	for _, hook := range evt.Hooks {
		go s.triggerHook(hook, body, headers)
	}
}

func (s *Server) triggerHook(hook types.Hook, body map[string]any, headers http.Header) {
	templateStr := s.config.GetTemplate(hook.Body)

	payload, err := render.ToMap(templateStr, body)
	if err != nil {
		log.Printf("[Render] Failed to parse rendered template to map (%s): %v", hook.Name, err)
		return
	}

	url := headers.Get(hook.EndpointKey)
	if url == "" {
		log.Printf("[Webhook] URL not found in header for key: %s", hook.EndpointKey)
		return
	}

	log.Printf("[Webhook] Triggering: %s", hook.Name)
	if err = postJSON(url, payload); err != nil {
		log.Printf("[Webhook] Failed to send request (%s): %v", hook.Name, err)
	}
}

func extractEventType(tmpl types.Template, headers http.Header, body map[string]any) (string, error) {
	switch tmpl.EventTypeIn {
	case "header":
		eventType := headers.Get(tmpl.EventTypeKey)
		if eventType == "" {
			return "", fmt.Errorf("event key '%s' not found in headers", tmpl.EventTypeKey)
		}
		return eventType, nil
	case "body":
		rawEvent, ok := body[tmpl.EventTypeKey]
		if !ok {
			return "", fmt.Errorf("event key '%s' not found in body", tmpl.EventTypeKey)
		}
		eventType, ok := rawEvent.(string)
		if !ok {
			return "", fmt.Errorf("event key '%s' is not a string in body", tmpl.EventTypeKey)
		}
		return eventType, nil
	default:
		return "", fmt.Errorf("unknown EventTypeIn value: '%s'", tmpl.EventTypeIn)
	}
}

func postJSON(url string, data map[string]any) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	return err
}
