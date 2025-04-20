package types

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Config []TemplateConfig

func (c Config) GetTemplates(receiver string, r *http.Request, payload []byte) (templates []TemplateConfig) {
	for _, template := range c {
		if template.Receiver != receiver {
			continue
		}

		if !isAuthorized(template.Auth, r, payload) {
			log.Printf("[AUTH] failed for receiver: %s with flow: %s", receiver, template.Auth.Flow)
			continue
		}

		templates = append(templates, template)
	}
	return
}

func isAuthorized(auth Auth, r *http.Request, payload []byte) bool {
	switch auth.Flow {
	case "none":
		return true

	case "basic auth":
		username, password, ok := r.BasicAuth()
		return ok && auth.Secret == fmt.Sprintf("%s:%s", username, password)

	case "gitlab":
		expected := sha512.Sum512([]byte(auth.Secret))
		actual := sha512.Sum512([]byte(r.Header.Get(auth.HeaderSecretKey)))
		return subtle.ConstantTimeCompare(actual[:], expected[:]) == 1

	case "github":
		signature := r.Header.Get(auth.HeaderSecretKey)
		if signature == "" {
			return false
		}
		signature = strings.TrimPrefix(signature, "sha256=")

		mac := hmac.New(sha256.New, []byte(auth.Secret))
		_, _ = mac.Write(payload)
		expectedMAC := hex.EncodeToString(mac.Sum(nil))

		return hmac.Equal([]byte(signature), []byte(expectedMAC))

	case "plain secret":
		return auth.Secret == r.Header.Get(auth.HeaderSecretKey)

	default:
		return false
	}
}
