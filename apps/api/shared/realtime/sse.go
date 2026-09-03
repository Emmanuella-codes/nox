package realtime

import (
	"bufio"
	"encoding/json"
	"fmt"
)

// writes one server-sent event frame to the stream.
func WriteEvent(w *bufio.Writer, event Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	if event.ID != "" {
		if _, err := fmt.Fprintf(w, "id: %s\n", event.ID); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintf(w, "data: %s\n\n", payload); err != nil {
		return err
	}
	return w.Flush()
}

// writes one comment heartbeat frame to the stream.
func WriteComment(w *bufio.Writer, value string) error {
	if _, err := fmt.Fprintf(w, ": %s\n\n", value); err != nil {
		return err
	}
	return w.Flush()
}
