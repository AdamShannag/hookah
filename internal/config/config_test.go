package config_test

import (
	"bytes"
	"github.com/AdamShannag/hookah/internal/auth"
	"github.com/AdamShannag/hookah/internal/config"
	"github.com/AdamShannag/hookah/internal/types"
	"net/http/httptest"
	"testing"
)

func TestNewConfig(t *testing.T) {
	tmpls := []types.Template{{Receiver: "discord", Auth: types.Auth{Flow: "none"}}}
	tmplMap := map[string]string{"discord": `{ "msg": "hello" }`}

	cfg := config.New(tmpls, tmplMap, auth.NewDefault())

	if len(cfg.GetConfigTemplates("discord", httptest.NewRequest("POST", "/", nil), nil)) != 1 {
		t.Error("expected one template to match")
	}
}

func TestGetTemplate(t *testing.T) {
	cfg := config.New(nil, map[string]string{
		"discord": `{ "msg": "hello" }`,
	}, auth.NewDefault())

	val := cfg.GetTemplate("discord")
	if val != `{ "msg": "hello" }` {
		t.Errorf("expected template body, got %s", val)
	}

	val = cfg.GetTemplate("unknown")
	if val != "{}" {
		t.Errorf("expected fallback '{}', got %s", val)
	}
}

func TestGetConfigTemplates_AuthFilter(t *testing.T) {
	templates := []types.Template{
		{Receiver: "slack", Auth: types.Auth{Flow: "plain secret"}},
		{Receiver: "slack", Auth: types.Auth{Flow: "none"}},
	}

	req := httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte("test")))

	t.Run("auth passes", func(t *testing.T) {
		cfg := config.New(templates, nil, auth.NewDefault())
		result := cfg.GetConfigTemplates("slack", req, []byte("test"))

		if len(result) != 2 {
			t.Errorf("expected 2 templates to pass auth, got %d", len(result))
		}
	})

	t.Run("auth fails", func(t *testing.T) {
		cfg := config.New([]types.Template{
			{Receiver: "slack", Auth: types.Auth{Flow: "no flow"}},
		}, nil, auth.NewDefault())

		result := cfg.GetConfigTemplates("slack", req, []byte("test"))

		if len(result) != 0 {
			t.Errorf("expected 0 templates to pass auth, got %d", len(result))
		}
	})

	t.Run("receiver mismatch", func(t *testing.T) {
		cfg := config.New(templates, nil, auth.NewDefault())
		result := cfg.GetConfigTemplates("discord", req, []byte("test"))

		if len(result) != 0 {
			t.Errorf("expected 0 templates for mismatched receiver, got %d", len(result))
		}
	})
}
