package blockchain

import "encoding/json"

// Message gets converted to/from JSON and sent in the body of pubsub messages.
type Message struct {
	Message  string
	SenderID string

	Annotations map[string]interface{}
}

type MessageOption func(cfg *Message) error

// Apply applies the given options to the config, returning the first error
// encountered (if any).
func (m *Message) Apply(opts ...MessageOption) error {
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(m); err != nil {
			return err
		}
	}
	return nil
}

func NewMessage(s string) *Message {
	return &Message{Message: s}
}

func (m *Message) Copy() *Message {
	copy := *m
	return &copy
}

func (m *Message) WithMessage(s string) *Message {
	copy := m.Copy()
	copy.Message = s
	return copy
}

func (m *Message) AnnotationsToObj(v interface{}) error {
	blob, err := json.Marshal(m.Annotations)
	if err != nil {
		return err
	}
	return json.Unmarshal(blob, v)
}
