package runtime

import (
	"context"
	"time"
)

func RunLoop(ctx context.Context, interval time.Duration, tick func(context.Context) error) error {
	if err := tick(ctx); err != nil {
		return err
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := tick(ctx); err != nil {
				return err
			}
		}
	}
}
