package condition

import (
	"net/http"
	"testing"
)

func TestEvaluate_SuccessCases(t *testing.T) {
	body := map[string]any{
		"objectMeta": map[string]any{
			"Labels": []any{
				map[string]any{"title": "API"},
				map[string]any{"title": "Backend"},
			},
		},
		"user": map[string]any{
			"name": "Alice",
		},
	}

	headers := http.Header{}
	headers.Set("x-gitlab-label", "API")
	headers.Set("x-user", "Alice")

	tests := []struct {
		name       string
		conditions []string
		want       bool
	}{
		{
			name:       "Header equals",
			conditions: []string{"{Header.x-user} {eq} {Body.user.name}"},
			want:       true,
		},
		{
			name:       "Header in label titles",
			conditions: []string{"{Header.x-gitlab-label} {in} {Body.objectMeta.Labels[].title}"},
			want:       true,
		},
		{
			name:       "Not equals",
			conditions: []string{"{Header.x-gitlab-label} {ne} {Body.user.name}"},
			want:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Evaluate(tt.conditions, headers, body)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestEvaluate_FailureCases(t *testing.T) {
	body := map[string]any{
		"objectMeta": map[string]any{
			"Labels": []any{
				map[string]any{"title": "API"},
			},
		},
	}
	headers := http.Header{}
	headers.Set("x-gitlab-label", "API")

	tests := []struct {
		name       string
		conditions []string
	}{
		{
			name:       "Invalid format",
			conditions: []string{"{Header.x-gitlab-label} {eq}"},
		},
		{
			name:       "Unsupported operator",
			conditions: []string{"{Header.x-gitlab-label} {gt} {Body.objectMeta.Labels}"},
		},
		{
			name:       "Invalid projection target",
			conditions: []string{"{Header.x-gitlab-label} {in} {Body.objectMeta.Invalid[].title}"},
		},
		{
			name:       "Right side not a list",
			conditions: []string{"{Header.x-gitlab-label} {in} {Body.objectMeta.Labels}"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, err := Evaluate(tt.conditions, headers, body)
			if err == nil && ok {
				t.Errorf("expected error but got success (ok = %v)", ok)
			}
		})
	}
}
