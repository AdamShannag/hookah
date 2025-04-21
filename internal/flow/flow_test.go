package flow_test

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/AdamShannag/hookah/internal/flow"
	"github.com/AdamShannag/hookah/internal/types"
	"net/http/httptest"
	"testing"
)

func TestNone(t *testing.T) {
	auth := types.Auth{}
	req := httptest.NewRequest("POST", "/", nil)
	if !flow.None(auth, req, nil) {
		t.Error("None should always return true")
	}
}

func TestBasicAuth(t *testing.T) {
	auth := types.Auth{Secret: "user:pass"}
	req := httptest.NewRequest("POST", "/", nil)
	req.SetBasicAuth("user", "pass")

	if !flow.BasicAuth(auth, req, nil) {
		t.Error("BasicAuth should succeed with matching credentials")
	}

	req.SetBasicAuth("user", "wrong")
	if flow.BasicAuth(auth, req, nil) {
		t.Error("BasicAuth should fail with incorrect credentials")
	}
}

func TestPlainSecret(t *testing.T) {
	auth := types.Auth{
		Secret:          "my-secret",
		HeaderSecretKey: "X-Secret-Key",
	}
	req := httptest.NewRequest("POST", "/", nil)
	req.Header.Set("X-Secret-Key", "my-secret")

	if !flow.PlainSecret(auth, req, nil) {
		t.Error("PlainSecret should succeed with correct secret")
	}

	req.Header.Set("X-Secret-Key", "wrong-secret")
	if flow.PlainSecret(auth, req, nil) {
		t.Error("PlainSecret should fail with incorrect secret")
	}
}

func TestGitlab(t *testing.T) {
	secret := "top-secret"
	auth := types.Auth{
		Secret:          secret,
		HeaderSecretKey: "X-Gitlab-Token",
	}
	req := httptest.NewRequest("POST", "/", nil)
	req.Header.Set("X-Gitlab-Token", secret)

	if !flow.Gitlab(auth, req, nil) {
		t.Error("Gitlab should succeed with correct hash")
	}

	req.Header.Set("X-Gitlab-Token", "invalid")
	if flow.Gitlab(auth, req, nil) {
		t.Error("Gitlab should fail with incorrect hash")
	}
}

func TestGithub(t *testing.T) {
	secret := "github-secret"
	auth := types.Auth{
		Secret:          secret,
		HeaderSecretKey: "X-Hub-Signature-256",
	}
	payload := []byte(`{"key":"value"}`)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	signature := hex.EncodeToString(mac.Sum(nil))
	req := httptest.NewRequest("POST", "/", bytes.NewReader(payload))
	req.Header.Set("X-Hub-Signature-256", "sha256="+signature)

	if !flow.Github(auth, req, payload) {
		t.Error("Github should succeed with correct signature")
	}

	req.Header.Set("X-Hub-Signature-256", "sha256=invalidsignature")
	if flow.Github(auth, req, payload) {
		t.Error("Github should fail with incorrect signature")
	}
}
