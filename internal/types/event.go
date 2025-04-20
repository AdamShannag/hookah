package types

type Events []Event

type Event struct {
	Event      string   `json:"event,omitempty"`
	Conditions []string `json:"conditions,omitempty"`
	Hooks      []Hook   `json:"hooks,omitempty"`
}

func (e Events) GetEvents(event string) (events []Event) {
	for _, evt := range e {
		if evt.Event == event {
			events = append(events, evt)
		}
	}

	return
}
