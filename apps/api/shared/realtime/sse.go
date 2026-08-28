package realtime

import (
	"bufio"
	"fmt"
)

// writes one server-sent event frame to the stream.
func WriteSSE(w *bufio.Writer, payload []byte) error {
	if _, err := fmt.Fprintf(w, "data: %s\n\n", payload); err != nil {
		return err
	}
	return w.Flush()
}

// writes one comment heartbeat frame to the stream.
func WriteSSEComment(w *bufio.Writer, value string) error {
	if _, err := fmt.Fprintf(w, ": %s\n\n", value); err != nil {
		return err
	}
	return w.Flush()
}
