package resolver

import (
	"reflect"
	"testing"
)

func TestPathResolver_Resolve(t *testing.T) {
	resolver := NewPathResolver()

	data := map[string]any{
		"user": map[string]any{
			"name": "Alice",
			"age":  30,
			"address": map[string]any{
				"city":  "Wonderland",
				"state": "Imagination",
			},
			"tags": []any{"admin", "user"},
			"friends": []any{
				map[string]any{"name": "Bob"},
				map[string]any{"name": "Charlie"},
			},
			"projects": []any{
				map[string]any{
					"name":  "ProjectX",
					"tasks": []any{"Design", "Develop", "Test"},
				},
				map[string]any{
					"name":  "ProjectY",
					"tasks": []any{"Plan", "Execute"},
				},
			},
		},
		"meta": map[string]any{
			"active": true,
		},
	}

	tests := []struct {
		name    string
		path    string
		want    any
		wantErr bool
	}{
		// Basic and nested access
		{"Simple value", "user.name", "Alice", false},
		{"Nested value", "user.address.city", "Wonderland", false},
		{"Missing key", "user.nonexistent", nil, true},
		{"Meta boolean", "meta.active", true, false},

		// Indexed array access
		{"Array access", "user.tags[1]", "user", false},
		{"Array out of bounds", "user.tags[5]", nil, true},

		// Projections
		{"Array projection", "user.friends[].name", []any{"Bob", "Charlie"}, false},
		{"Projection with nested array", "user.projects[].tasks[0]", []any{"Design", "Plan"}, false},
		{"Projection full task lists", "user.projects[].tasks", []any{"Design", "Develop", "Test", "Plan", "Execute"}, false},

		// Mixed nested indexing
		{"Deep access with index", "user.projects[1].tasks[1]", "Execute", false},
		{"Invalid nested array", "user.projects[1].unknown[0]", nil, true},

		// Invalid path formats
		{"Non-array projection", "user.name[]", nil, true},
		{"Index on non-array", "user.name[0]", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolver.Resolve(tt.path, data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Resolve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Resolve() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
