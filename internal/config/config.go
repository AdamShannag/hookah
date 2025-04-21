package config

import (
	"github.com/AdamShannag/hookah/internal/auth"
	"github.com/AdamShannag/hookah/internal/types"
	"log"
	"net/http"
)

type Config struct {
	templateConfigs []types.Template
	templates       map[string]string
	auth            auth.Auth
}

func New(templateConfigs []types.Template, templates map[string]string, auth auth.Auth) *Config {
	return &Config{
		templateConfigs: templateConfigs,
		templates:       templates,
		auth:            auth,
	}
}

func (c *Config) GetTemplate(template string) string {
	body, ok := c.templates[template]
	if !ok {
		return "{}"
	}
	return body
}

func (c *Config) GetConfigTemplates(receiver string, r *http.Request, payload []byte) (templates []types.Template) {
	for _, template := range c.templateConfigs {
		if template.Receiver != receiver {
			continue
		}

		if !c.auth.ApplyFlow(template.Auth, r, payload) {
			log.Printf("[AUTH] failed for receiver: %s with flow: %s", receiver, template.Auth.Flow)
			continue
		}

		templates = append(templates, template)
	}
	return
}
