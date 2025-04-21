package render

import (
	"log"
	"strings"
	"text/template"
	"time"
)

var funcMap = template.FuncMap{
	"now":       now,
	"format":    format,
	"parseTime": parseTime,
	"pastTense": pastTense,
	"lower":     strings.ToLower,
	"upper":     strings.ToUpper,
	"title":     strings.ToTitle,
	"trim":      strings.TrimSpace,
	"contains":  strings.Contains,
	"replace":   strings.ReplaceAll,
	"default":   defaultValue,
}

func now() time.Time { return time.Now() }

func format(t time.Time, format string) string {
	return t.Format(format)
}

func parseTime(tm string, layout string) time.Time {
	t, err := time.Parse(layout, tm)
	if err != nil {
		log.Printf("[Template] %v", err)
		return time.Time{}
	}

	return t
}

func pastTense(word string) string {
	if len(word) > 0 && word[len(word)-1] == 'e' {
		return word + "d"
	}

	return word + "ed"
}

func defaultValue(val, fallback any) any {
	if val == nil || val == "" {
		return fallback
	}
	return val
}
