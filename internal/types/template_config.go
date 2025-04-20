package types

type TemplateConfig struct {
	Receiver     string `json:"receiver"`
	Auth         Auth   `json:"auth"`
	EventTypeIn  string `json:"event_type_in"`
	EventTypeKey string `json:"event_type_key"`
	Events       Events `json:"events,omitempty"`
}

type Hook struct {
	Name        string         `json:"name"`
	EndpointKey string         `json:"endpoint_key"`
	Body        map[string]any `json:"body,omitempty"`
}
