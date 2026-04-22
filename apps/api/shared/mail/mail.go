package mail

import "context"

type Message struct {
	ToEmail     string
	ToName      string
	Subject     string
	HTMLContent string
	TextContent string
}

type Provider interface {
	Send(ctx context.Context, message Message) error
}
