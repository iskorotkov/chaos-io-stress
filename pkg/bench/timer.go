package bench

import (
	"time"
)

func Timer(d time.Duration, done <-chan struct{}) <-chan struct{} {
	res := make(chan struct{})

	go func() {
		defer close(res)

		time.Sleep(d)

		for {
			select {
			case <-done:
				return
			case res <- struct{}{}:
				time.Sleep(d)
			}
		}
	}()

	return res
}
