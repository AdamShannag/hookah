package auth

import (
	"github.com/AdamShannag/hookah/internal/flow"
	"github.com/AdamShannag/hookah/internal/types"
	"net/http"
)

type Auth interface {
	RegisterFlow(flow string, flowFunc flow.Func) Auth
	ApplyFlow(auth types.Auth, r *http.Request, payload []byte) bool
}

type auth struct {
	flows map[string]flow.Func
}

func New() Auth {
	return &auth{make(map[string]flow.Func)}
}

func NewDefault() Auth {
	return New().
		RegisterFlow("none", flow.None).
		RegisterFlow("plain secret", flow.PlainSecret).
		RegisterFlow("basic auth", flow.BasicAuth).
		RegisterFlow("gitlab", flow.Gitlab).
		RegisterFlow("github", flow.Github)
}

func (a *auth) RegisterFlow(flow string, flowFunc flow.Func) Auth {
	a.flows[flow] = flowFunc
	return a
}

func (a *auth) ApplyFlow(auth types.Auth, r *http.Request, payload []byte) bool {
	flowFunc, ok := a.flows[auth.Flow]
	if !ok {

		return false
	}

	return flowFunc(auth, r, payload)
}
