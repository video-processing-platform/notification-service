package mail

import (
	"context"
	"time"
)

func retry(
	ctx context.Context,
	attempts int,
	delay time.Duration,
	fn func() error,
) error {

	var err error

	for i := 0; i < attempts; i++ {

		if ctx.Err() != nil {
			return ctx.Err()
		}

		err = fn()
		if err == nil {
			return nil
		}

		if i == attempts-1 {
			break
		}

		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return err
}
