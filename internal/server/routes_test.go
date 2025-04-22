package server

import (
	"bytes"
	"encoding/json"
	"github.com/AdamShannag/hookah/internal/auth"
	"github.com/AdamShannag/hookah/internal/condition"
	"github.com/AdamShannag/hookah/internal/config"
	"github.com/AdamShannag/hookah/internal/resolver"
	"github.com/AdamShannag/hookah/internal/types"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
	"time"
)

func TestWebhookHandler_DispatchesToSimulatedEndpoint(t *testing.T) {
	var (
		receivedPayload map[string]any
		mu              sync.Mutex
		wg              sync.WaitGroup
	)

	wg.Add(1)
	mockDiscord := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, _ := io.ReadAll(r.Body)

		mu.Lock()
		defer mu.Unlock()
		_ = json.Unmarshal(body, &receivedPayload)
		wg.Done()
	}))
	defer mockDiscord.Close()

	testServer := &Server{
		evaluator: condition.NewDefaultEvaluator(resolver.NewPathResolver()),
		config: config.New([]types.Template{
			{
				Receiver:     "gitlab",
				Auth:         types.Auth{Flow: "none"},
				EventTypeKey: "event_name",
				EventTypeIn:  "body",
				Events: types.Events{
					{
						Event:      "issue",
						Conditions: []string{"{Header.X-Custom} {eq} {active}"},
						Hooks: []types.Hook{
							{
								Name:        "MockDiscord",
								EndpointKey: "Webhook-URL",
								Body:        "discord.tmpl",
							},
						},
					},
				},
			},
		}, map[string]string{
			"discord.tmpl": getBodyTemplate("Issue received"),
		}, auth.NewDefault()),
	}

	reqBody := map[string]any{
		"event_name": "issue",
		"status":     "active",
	}
	bodyBuf := new(bytes.Buffer)
	_ = json.NewEncoder(bodyBuf).Encode(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/webhooks/gitlab", bodyBuf)
	req.Header.Set("X-Custom", "active")
	req.Header.Set("Webhook-URL", mockDiscord.URL)

	rr := httptest.NewRecorder()
	testServer.RegisterRoutes().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	wg.Wait()

	mu.Lock()
	defer mu.Unlock()
	if receivedPayload["content"] != "Issue received" {
		t.Fatalf("expected payload 'Issue received', got: %v", receivedPayload)
	}
}

func TestWebhookHandler_DoesNotDispatchWhenConditionFails(t *testing.T) {
	var wasCalled bool
	wg := sync.WaitGroup{}

	mockDiscord := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wasCalled = true
		wg.Done()
	}))
	defer mockDiscord.Close()

	testServer := &Server{
		evaluator: condition.NewDefaultEvaluator(resolver.NewPathResolver()),
		config: config.New([]types.Template{
			{
				Receiver:     "gitlab",
				Auth:         types.Auth{Flow: "none"},
				EventTypeKey: "event_name",
				EventTypeIn:  "body",
				Events: types.Events{
					{
						Event:      "issue",
						Conditions: []string{"{Header.X-Custom} {eq} {Body.status}"},
						Hooks: []types.Hook{
							{
								Name:        "MockDiscord",
								EndpointKey: "Webhook-URL",
								Body:        "discord.tmpl",
							},
						},
					},
				},
			},
		}, map[string]string{
			"discord.tmpl": getBodyTemplate("Should not be triggered"),
		}, auth.NewDefault()),
	}

	reqBody := map[string]any{
		"event_name": "issue",
		"status":     "inactive",
	}
	bodyBuf := new(bytes.Buffer)
	_ = json.NewEncoder(bodyBuf).Encode(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/webhooks/gitlab", bodyBuf)
	req.Header.Set("X-Custom", "active")
	req.Header.Set("Webhook-URL", mockDiscord.URL)

	rr := httptest.NewRecorder()
	testServer.RegisterRoutes().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	time.Sleep(100 * time.Millisecond)

	if wasCalled {
		t.Fatalf("expected no webhook call, but it was triggered")
	}
}

func TestWebhookHandler_UsesQueryParamsAsHeaders(t *testing.T) {
	var (
		receivedPayload map[string]any
		mu              sync.Mutex
		wg              sync.WaitGroup
	)

	wg.Add(1)
	mockDiscord := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, _ := io.ReadAll(r.Body)

		mu.Lock()
		defer mu.Unlock()
		_ = json.Unmarshal(body, &receivedPayload)
		wg.Done()
	}))
	defer mockDiscord.Close()

	testServer := &Server{
		evaluator: condition.NewDefaultEvaluator(resolver.NewPathResolver()),
		config: config.New([]types.Template{
			{
				Receiver:     "gitlab",
				Auth:         types.Auth{Flow: "none"},
				EventTypeKey: "event_name",
				EventTypeIn:  "body",
				Events: types.Events{
					{
						Event:      "issue",
						Conditions: []string{"{Header.X-Custom} {eq} {active}"},
						Hooks: []types.Hook{
							{
								Name:        "MockDiscord",
								EndpointKey: "Webhook-URL",
								Body:        "discord.tmpl",
							},
						},
					},
				},
			},
		}, map[string]string{
			"discord.tmpl": getBodyTemplate("Query param test passed"),
		}, auth.NewDefault()),
	}

	reqBody := map[string]any{
		"event_name": "issue",
	}
	bodyBuf := new(bytes.Buffer)
	_ = json.NewEncoder(bodyBuf).Encode(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/webhooks/gitlab?X-Custom=active&Webhook-URL="+url.QueryEscape(mockDiscord.URL), bodyBuf)

	rr := httptest.NewRecorder()
	testServer.RegisterRoutes().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	wg.Wait()

	mu.Lock()
	defer mu.Unlock()
	if receivedPayload["content"] != "Query param test passed" {
		t.Fatalf("expected payload 'Query param test passed', got: %v", receivedPayload)
	}
}

func getBodyTemplate(content string) string {
	marshal, _ := json.Marshal(map[string]string{
		"content": content,
	})

	return string(marshal)
}
