package shared

import (
	"encoding/json"
	"testing"
)

func TestPipeResSerializesLowercaseContractKeys(t *testing.T) {
	data := map[string]string{"id": "user-1"}
	res := PipeRes[map[string]string]{
		Success: true,
		Message: CreatePipeMessage("ok"),
		Data:    &data,
	}

	body, err := json.Marshal(res)
	if err != nil {
		t.Fatalf("marshal PipeRes: %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(body, &decoded); err != nil {
		t.Fatalf("unmarshal PipeRes: %v", err)
	}

	for _, key := range []string{"success", "message", "data"} {
		if _, ok := decoded[key]; !ok {
			t.Fatalf("expected JSON key %q in %s", key, string(body))
		}
	}

	for _, key := range []string{"Success", "Message", "Data"} {
		if _, ok := decoded[key]; ok {
			t.Fatalf("did not expect capitalized JSON key %q in %s", key, string(body))
		}
	}
}
