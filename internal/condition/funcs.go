package condition

import (
	"fmt"
	"strings"
)

func equals(a, b any) (bool, error) {
	return a == b, nil
}

func notEquals(a, b any) (bool, error) {
	return a != b, nil
}

func in(a any, b any) (bool, error) {
	list, ok := b.([]any)
	if !ok {
		return false, fmt.Errorf("right side must be a list")
	}
	for _, item := range list {
		if item == a {
			return true, nil
		}
	}
	return false, nil
}

func contains(a, b any) (bool, error) {
	as, aok := a.(string)
	bs, bok := b.(string)
	if !aok || !bok {
		return false, fmt.Errorf("both sides must be strings")
	}
	return strings.Contains(as, bs), nil
}

func notIn(a any, b any) (bool, error) {
	list, ok := b.([]any)
	if !ok {
		return false, fmt.Errorf("right side must be a list")
	}
	for _, item := range list {
		if item == a {
			return false, nil
		}
	}
	return true, nil
}

func startsWith(a, b any) (bool, error) {
	as, aok := a.(string)
	bs, bok := b.(string)
	if !aok || !bok {
		return false, fmt.Errorf("both sides must be strings")
	}
	return strings.HasPrefix(as, bs), nil
}

func endsWith(a, b any) (bool, error) {
	as, aok := a.(string)
	bs, bok := b.(string)
	if !aok || !bok {
		return false, fmt.Errorf("both sides must be strings")
	}
	return strings.HasSuffix(as, bs), nil
}
