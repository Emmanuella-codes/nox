package shared

import "github.com/rs/zerolog/log"

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

func PipeSuccess[T any](message PipeMessage, data *T) *PipeRes[T] {
	return &PipeRes[T]{Success: true, Message: message, Data: data}
}

func PipeError[T any](message PipeMessage) *PipeRes[T] {
	return &PipeRes[T]{Success: false, Message: message}
}

func PipeInternalError[T any](err error, domain string, operation string, message PipeMessage) *PipeRes[T] {
	if err != nil {
		log.Error().Err(err).Str("operation", operation).Msg(domain + " internal error")
	}
	return PipeError[T](message)
}
