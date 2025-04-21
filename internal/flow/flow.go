package flow

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"github.com/AdamShannag/hookah/internal/types"
	"net/http"
	"strings"
)

type Func func(auth types.Auth, r *http.Request, payload []byte) bool

func None(_ types.Auth, _ *http.Request, _ []byte) bool {
	return true
}

func BasicAuth(auth types.Auth, r *http.Request, _ []byte) bool {
	username, password, ok := r.BasicAuth()
	return ok && auth.Secret == fmt.Sprintf("%s:%s", username, password)
}

func PlainSecret(auth types.Auth, r *http.Request, _ []byte) bool {
	return auth.Secret == r.Header.Get(auth.HeaderSecretKey)
}

func Gitlab(auth types.Auth, r *http.Request, _ []byte) bool {
	expected := sha512.Sum512([]byte(auth.Secret))
	actual := sha512.Sum512([]byte(r.Header.Get(auth.HeaderSecretKey)))
	return subtle.ConstantTimeCompare(actual[:], expected[:]) == 1
}

func Github(auth types.Auth, r *http.Request, payload []byte) bool {
	signature := r.Header.Get(auth.HeaderSecretKey)
	if signature == "" {
		return false
	}
	signature = strings.TrimPrefix(signature, "sha256=")

	mac := hmac.New(sha256.New, []byte(auth.Secret))
	_, _ = mac.Write(payload)
	expectedMAC := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedMAC))
}
