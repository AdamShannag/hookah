package auth_test

import (
	"bytes"
	"github.com/AdamShannag/hookah/internal/auth"
	"github.com/AdamShannag/hookah/internal/types"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterAndApplyFlow(t *testing.T) {
	a := auth.New()

	a.RegisterFlow("mock-flow-pass", mockFlow(true))
	a.RegisterFlow("mock-flow-fail", mockFlow(false))

	req := httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte(`{"data":"test"}`)))

	tests := []struct {
		name     string
		auth     types.Auth
		expected bool
	}{
		{"FlowPasses", types.Auth{Flow: "mock-flow-pass"}, true},
		{"FlowFails", types.Auth{Flow: "mock-flow-fail"}, false},
		{"FlowMissing", types.Auth{Flow: "unregistered"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := a.ApplyFlow(tt.auth, req, []byte("test"))
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func mockFlow(expected bool) func(types.Auth, *http.Request, []byte) bool {
	return func(a types.Auth, r *http.Request, b []byte) bool {
		return expected
	}
}
