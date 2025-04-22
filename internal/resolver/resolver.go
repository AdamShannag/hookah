package resolver

import (
	"fmt"
	"strconv"
	"strings"
)

// Resolver defines the interface for any value resolver.
type Resolver interface {
	Resolve(path string, data any) (any, error)
}

// pathResolver resolves dotted paths and supports array access and projection.
type pathResolver struct{}

func NewPathResolver() Resolver {
	return &pathResolver{}
}

// Resolve a value by navigating the provided path (e.g. "users[0].name" or "users[].name").
func (r *pathResolver) Resolve(path string, data any) (any, error) {
	parts := strings.Split(path, ".")
	return resolvePath(data, parts)
}

// resolvePath resolves the nested path on a given data structure.
func resolvePath(current any, parts []string) (any, error) {
	for i, part := range parts {
		key, index, isIndexed := parseIndexedKey(part)

		objMap, ok := current.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("expected map at '%s'", part)
		}

		value, exists := objMap[key]
		if !exists {
			return nil, fmt.Errorf("key '%s' not found", key)
		}

		if !isIndexed {
			current = value
			continue
		}

		array, ok := value.([]any)
		if !ok {
			return nil, fmt.Errorf("expected array at '%s'", key)
		}

		if index == nil {
			return projectArray(array, parts[i+1:])
		}

		if *index < 0 || *index >= len(array) {
			return nil, fmt.Errorf("index out of bounds at '%s[%d]'", key, *index)
		}

		current = array[*index]
	}

	return current, nil
}

// projectArray handles projection through an array of maps using the remaining path parts.
func projectArray(array []any, remainingParts []string) ([]any, error) {
	results := make([]any, 0, len(array))
	for _, item := range array {
		val, err := resolvePath(item, remainingParts)
		if err != nil {
			continue
		}
		if list, ok := val.([]any); ok {
			results = append(results, list...)
		} else {
			results = append(results, val)
		}
	}
	return results, nil
}

// parseIndexedKey parses keys like "users[0]" or "users[]" and returns the base key, index (if any), and a flag.
func parseIndexedKey(part string) (key string, index *int, isIndexed bool) {
	start := strings.Index(part, "[")
	end := strings.Index(part, "]")

	// No brackets
	if start == -1 || end == -1 || end < start {
		return part, nil, false
	}

	key = part[:start]
	idxStr := part[start+1 : end]

	if idxStr == "" {
		return key, nil, true
	}

	i, err := strconv.Atoi(idxStr)
	if err != nil {
		return part, nil, false
	}

	return key, &i, true
}
