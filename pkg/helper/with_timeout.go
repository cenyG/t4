package helper

import (
	"context"
	"time"
)

func WithTimeout[T any](ctx context.Context, timeout time.Duration, fn func() (T, error)) (T, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)

	errCh := make(chan error, 1)
	defer close(errCh)
	resCh := make(chan T, 1)
	defer close(resCh)

	go func() {
		res, err := fn()
		if err != nil {
			errCh <- err
			return
		}
		select {
		case <-ctx.Done():
			errCh <- ctx.Err()
		case resCh <- res:
		}
	}()

	var nilValue T
	select {
	case err := <-errCh:
		cancel()
		return nilValue, err
	case r := <-resCh:
		cancel()
		return r, nil
	}
}
