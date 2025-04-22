package condition_test

import (
	"github.com/AdamShannag/hookah/internal/condition"
	"github.com/AdamShannag/hookah/internal/resolver"
	"net/http"
	"testing"
)

func TestEvaluator_EvaluateAll(t *testing.T) {
	res := resolver.NewPathResolver()
	eval := condition.NewDefaultEvaluator(res)

	tests := []struct {
		name       string
		conditions []string
		headers    http.Header
		body       map[string]any
		want       bool
		wantErr    bool
	}{
		{
			name:       "Equal condition true",
			conditions: []string{"{Body.user.name} {eq} {Jane}"},
			body: map[string]any{
				"user": map[string]any{"name": "Jane"},
			},
			want: true,
		},
		{
			name:       "Equal condition false",
			conditions: []string{"{Body.user.name} {eq} {Bob}"},
			body: map[string]any{
				"user": map[string]any{"name": "Alice"},
			},
			want: false,
		},
		{
			name:       "Not equal true",
			conditions: []string{"{Body.user.name} {ne} {Bob}"},
			body: map[string]any{
				"user": map[string]any{"name": "Alice"},
			},
			want: true,
		},
		{
			name:       "Header value match",
			conditions: []string{"{Header.X-Token} {eq} {abc123}"},
			headers:    http.Header{"X-Token": []string{"abc123"}},
			want:       true,
		},
		{
			name:       "Unsupported operator",
			conditions: []string{"{Body.user.name} {gt} {Jane}"},
			body:       map[string]any{"user": map[string]any{"name": "Jane"}},
			wantErr:    true,
		},
		{
			name:       "Malformed condition",
			conditions: []string{"{Body.user.name} [eq] {Jane}"},
			body:       map[string]any{"user": map[string]any{"name": "Jane"}},
			wantErr:    true,
		},
		{
			name:       "Missing field",
			conditions: []string{"{Body.user.age} {eq} {30}"},
			body:       map[string]any{"user": map[string]any{}},
			want:       false,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := eval.EvaluateAll(tt.conditions, tt.headers, tt.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("EvaluateAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("EvaluateAll() = %v, want %v", got, tt.want)
			}
		})
	}
}
