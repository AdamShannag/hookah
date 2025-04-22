package condition

import (
	"fmt"
	"github.com/AdamShannag/hookah/internal/resolver"
	"net/http"
	"strings"
)

type Evaluator interface {
	Register(op string, fn OperatorFunc) Evaluator
	EvaluateAll(conditions []string, headers http.Header, body map[string]any) (bool, error)
}
type OperatorFunc func(left, right any) (bool, error)

type evaluator struct {
	operators map[string]OperatorFunc
	resolver  resolver.Resolver
}

func NewEvaluator(resolver resolver.Resolver) Evaluator {
	return &evaluator{operators: map[string]OperatorFunc{}, resolver: resolver}
}

func NewDefaultEvaluator(resolver resolver.Resolver) Evaluator {
	return NewEvaluator(resolver).
		Register("{eq}", equals).
		Register("{ne}", notEquals).
		Register("{in}", in).
		Register("{notIn}", notIn).
		Register("{contains}", contains).
		Register("{startsWith}", startsWith).
		Register("{endsWith}", endsWith)
}

func (e *evaluator) Register(op string, fn OperatorFunc) Evaluator {
	e.operators[op] = fn
	return e
}

func (e *evaluator) EvaluateAll(conditions []string, headers http.Header, body map[string]any) (bool, error) {
	for _, cond := range conditions {
		match, err := e.evaluateOne(cond, headers, body)
		if err != nil {
			return false, err
		}
		if !match {
			return false, nil
		}
	}
	return true, nil
}

func (e *evaluator) evaluateOne(condition string, headers http.Header, body map[string]any) (bool, error) {
	op, leftRaw, rightRaw, err := e.extractParts(condition)
	if err != nil {
		return false, err
	}

	leftVal, err := e.resolveValue(leftRaw, headers, body)
	if err != nil {
		return false, fmt.Errorf("left value: %w", err)
	}

	rightVal, err := e.resolveValue(rightRaw, headers, body)
	if err != nil {
		return false, fmt.Errorf("right value: %w", err)
	}

	fn, ok := e.operators[op]
	if !ok {
		return false, fmt.Errorf("unknown operator: %s", op)
	}

	return fn(leftVal, rightVal)
}

func (e *evaluator) extractParts(condition string) (op, left, right string, err error) {
	for opr := range e.operators {
		if strings.Contains(condition, opr) {
			parts := strings.SplitN(condition, opr, 2)
			if len(parts) != 2 {
				return "", "", "", fmt.Errorf("invalid condition: %s", condition)
			}
			return opr, strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), nil
		}
	}
	return "", "", "", fmt.Errorf("unsupported operator in: %s", condition)
}

func (e *evaluator) resolveValue(expr string, headers http.Header, body map[string]any) (any, error) {
	expr = strings.Trim(expr, "{}")
	switch {
	case strings.HasPrefix(expr, "Header."):
		key := strings.TrimPrefix(expr, "Header.")
		return headers.Get(key), nil
	case strings.HasPrefix(expr, "Body."):
		return e.resolver.Resolve(strings.TrimPrefix(expr, "Body."), body)
	default:
		return expr, nil
	}
}
