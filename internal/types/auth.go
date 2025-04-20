package types

type Auth struct {
	Flow            string `json:"flow"`
	HeaderSecretKey string `json:"header_secret_key,omitempty"`
	Secret          string `json:"secret"`
}
