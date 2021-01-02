package bench

import (
	"sync/atomic"
	"time"
)

func Benchmark(f func(), callback func(count int64), d time.Duration, done <-chan struct{}) {
	var ops int64 = 0

	stopTimer := make(chan struct{})
	defer close(stopTimer)
	timer := Timer(d, stopTimer)

	stopWorker := make(chan struct{})
	defer close(stopWorker)
	go startWorker(f, &ops, stopWorker)

	for {
		select {
		case <-done:
			return
		case <-timer:
			callback(atomic.SwapInt64(&ops, 0))
		}
	}
}

func startWorker(f func(), ops *int64, stopWorker chan struct{}) {
	for {
		select {
		case <-stopWorker:
			return
		default:
			f()
			atomic.AddInt64(ops, 1)
		}
	}
}
