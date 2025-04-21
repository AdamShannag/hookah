package types

type Template struct {
	Receiver     string `json:"receiver"`
	Auth         Auth   `json:"auth"`
	EventTypeIn  string `json:"event_type_in"`
	EventTypeKey string `json:"event_type_key"`
	Events       Events `json:"events,omitempty"`
}

type Hook struct {
	Name        string `json:"name"`
	EndpointKey string `json:"endpoint_key"`
	Body        string `json:"body,omitempty"`
}

type Auth struct {
	Flow            string `json:"flow"`
	HeaderSecretKey string `json:"header_secret_key,omitempty"`
	Secret          string `json:"secret"`
}
