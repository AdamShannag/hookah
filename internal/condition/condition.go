package condition

import (
	"fmt"
	"net/http"
	"strings"
)

// Supported condition operators
var supportedOps = []string{"{eq}", "{ne}", "{in}"}

// Evaluate checks if all provided conditions match
func Evaluate(conditions []string, headers http.Header, body map[string]any) (bool, error) {
	for _, condition := range conditions {
		ok, err := evaluate(condition, headers, body)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}
	return true, nil
}

// evaluate processes a single condition
func evaluate(condition string, headers http.Header, body map[string]any) (bool, error) {
	operator, leftExpr, rightExpr, err := splitCondition(condition)
	if err != nil {
		return false, err
	}

	leftVal, err := resolvePlaceholder(leftExpr, headers, body)
	if err != nil {
		return false, fmt.Errorf("resolving left side: %w", err)
	}

	rightVal, err := resolvePlaceholder(rightExpr, headers, body)
	if err != nil {
		return false, fmt.Errorf("resolving right side: %w", err)
	}

	return applyOperator(operator, leftVal, rightVal)
}

// splitCondition separates a condition into left, operator, right
func splitCondition(condition string) (string, string, string, error) {
	for _, op := range supportedOps {
		if strings.Contains(condition, op) {
			parts := strings.Split(condition, op)
			if len(parts) != 2 {
				return "", "", "", fmt.Errorf("invalid condition format: %s", condition)
			}
			return op, strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), nil
		}
	}
	return "", "", "", fmt.Errorf("unsupported operator in: %s", condition)
}

// applyOperator evaluates the actual condition
func applyOperator(op string, left, right any) (bool, error) {
	switch op {
	case "{eq}":
		return left == right, nil
	case "{ne}":
		return left != right, nil
	case "{in}":
		rightList, ok := right.([]any)
		if !ok {
			return false, fmt.Errorf("right side is not a list: %v", right)
		}
		for _, item := range rightList {
			if item == left {
				return true, nil
			}
		}
		return false, nil
	default:
		return false, fmt.Errorf("unknown operator: %s", op)
	}
}

// resolvePlaceholder extracts the actual value from headers or body
func resolvePlaceholder(placeholder string, headers http.Header, body map[string]any) (any, error) {
	expr := strings.Trim(placeholder, "{}")

	switch {
	case strings.HasPrefix(expr, "Header."):
		return resolveHeaderValue(expr, headers)
	case strings.HasPrefix(expr, "Body."):
		return resolveBodyValue(expr, body)
	default:
		return expr, nil
	}
}

func resolveHeaderValue(expr string, headers http.Header) (string, error) {
	key := strings.TrimPrefix(expr, "Header.")
	return headers.Get(key), nil
}

func resolveBodyValue(expr string, body map[string]any) (any, error) {
	path := strings.TrimPrefix(expr, "Body.")
	keys := strings.Split(path, ".")

	var current any = body

	for i := 0; i < len(keys); i++ {
		key := keys[i]

		// Handle array projection like labels[].title
		if strings.HasSuffix(key, "[]") {
			baseKey := strings.TrimSuffix(key, "[]")
			nextKey := keys[i+1]

			objMap, ok := current.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("expected map at %v", keys[:i])
			}
			arr, ok := objMap[baseKey].([]any)
			if !ok {
				return nil, fmt.Errorf("expected slice for key '%s'", baseKey)
			}

			var results []any
			for _, item := range arr {
				if m, ok := item.(map[string]any); ok {
					results = append(results, m[nextKey])
				}
			}
			return results, nil
		}

		// Regular nested field
		m, ok := current.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("expected map at path: %v", keys[:i])
		}
		current = m[key]
	}

	return current, nil
}
