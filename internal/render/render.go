package render

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"
)

func ToMap(tmplStr string, dataSource map[string]any) (map[string]any, error) {
	tmpl, err := template.New("map-template").Funcs(funcMap).Parse(tmplStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing template: %w", err)
	}

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, dataSource); err != nil {
		return nil, fmt.Errorf("error executing template: %w", err)
	}

	var result map[string]any
	if err = json.Unmarshal(buf.Bytes(), &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling rendered template: %w", err)
	}

	return result, nil
}
