package push

import (
	"context"
	"errors"

	"github.com/emmanuella-codes/nox/models"
)

var ErrInvalidToken = errors.New("push token invalid")

type Payload struct {
	Title      string
	Body       string
	TargetPath string
	Badge      int
	Raw        []byte
}

type Provider interface {
	// Send delivers one push payload to one device.
	Send(ctx context.Context, device *models.NotificationDevice, payload Payload) error
}
