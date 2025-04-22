package resolver

import (
	"testing"
)

var benchmarkData = map[string]any{
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
}

func BenchmarkResolveSimplePath(b *testing.B) {
	resolver := NewPathResolver()
	for i := 0; i < b.N; i++ {
		_, _ = resolver.Resolve("user.name", benchmarkData)
	}
}

func BenchmarkResolveNestedPath(b *testing.B) {
	resolver := NewPathResolver()
	for i := 0; i < b.N; i++ {
		_, _ = resolver.Resolve("user.address.city", benchmarkData)
	}
}

func BenchmarkResolveIndexedPath(b *testing.B) {
	resolver := NewPathResolver()
	for i := 0; i < b.N; i++ {
		_, _ = resolver.Resolve("user.projects[1].tasks[1]", benchmarkData)
	}
}

func BenchmarkResolveProjectedPath(b *testing.B) {
	resolver := NewPathResolver()
	for i := 0; i < b.N; i++ {
		_, _ = resolver.Resolve("user.projects[].tasks", benchmarkData)
	}
}
