package shared

type PipeMessage string

type PipeRes[T any] struct {
	Success  bool        `json:"success"`
	Message  PipeMessage `json:"message"`
	Data     *T          `json:"data,omitempty"`
	HookData any         `json:"hook_data,omitempty"`
	Token    string      `json:"token,omitempty"`
}

func CreatePipeMessage(msg string) PipeMessage {
	return PipeMessage(msg)
}
