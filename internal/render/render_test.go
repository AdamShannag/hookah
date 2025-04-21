package render

import (
	"strings"
	"testing"
)

func TestToMap_Success(t *testing.T) {
	tmpl := `{"msg": "hello {{.name}}", "active": {{.active}}}`
	data := map[string]any{
		"name":   "world",
		"active": true,
	}

	result, err := ToMap(tmpl, data)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if result["msg"] != "hello world" {
		t.Errorf("unexpected msg: %v", result["msg"])
	}
	if result["active"] != true {
		t.Errorf("unexpected active: %v", result["active"])
	}
}

func TestToMap_InvalidTemplate(t *testing.T) {
	_, err := ToMap(`{{.name`, map[string]any{"name": "fail"})
	if err == nil || !strings.Contains(err.Error(), "parsing template") {
		t.Fatalf("expected template parsing error, got: %v", err)
	}
}

func TestToMap_TemplateExecError(t *testing.T) {
	tmpl := `{{call .badFunc}}`
	data := map[string]any{}

	_, err := ToMap(tmpl, data)
	if err == nil {
		t.Fatal("expected template execution error, got nil")
	}
}

func TestToMap_InvalidJSON(t *testing.T) {
	tmpl := `{{.name}}`
	data := map[string]any{"name": "no-quotes"}

	_, err := ToMap(tmpl, data)
	if err == nil || !strings.Contains(err.Error(), "unmarshaling") {
		t.Fatalf("expected JSON unmarshal error, got: %v", err)
	}
}
