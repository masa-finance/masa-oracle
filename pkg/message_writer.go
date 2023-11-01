package masa

import "github.com/mudler/edgevpn/pkg/hub"

// messageWriter is a struct returned by the node that satisfies the io.Writer interface
// on the underlying hub.
// Everything Write into the message writer is enqueued to a message channel
// which is sealed and processed by the node
type messageWriter struct {
	input chan<- *hub.Message
	mess  *hub.Message
}

// Write writes a slice of bytes to the message channel
func (mw *messageWriter) Write(p []byte) (n int, err error) {
	return mw.Send(mw.mess.WithMessage(string(p)))
}

// Send sends a message to the channel
func (mw *messageWriter) Send(copy *hub.Message) (n int, err error) {
	mw.input <- copy
	return len(copy.Message), nil
}
