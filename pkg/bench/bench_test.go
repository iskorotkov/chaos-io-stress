package bench

import (
	"testing"
	"time"
)

func TestBenchmark(t *testing.T) {
	createTimer := func(d time.Duration) <-chan struct{} {
		ch := make(chan struct{})

		go func() {
			defer close(ch)

			<-time.After(d)
			ch <- struct{}{}
		}()

		return ch
	}

	type args struct {
		f        func()
		callback func(count int64)
		d        time.Duration
		done     <-chan struct{}
	}
	tests := []struct {
		name string
		args args
	}{
		{"single iteration", args{
			f: func() {
				time.Sleep(time.Millisecond * 10)
			},
			callback: func(count int64) {
				if count < 90 || count > 105 {
					t.Fatalf("expected %d-%d, actual %d", 90, 105, count)
				}
			},
			d:    time.Millisecond * 1000,
			done: createTimer(time.Millisecond * 1000),
		}},
		{"several iterations", args{
			f: func() {
				time.Sleep(time.Millisecond * 10)
			},
			callback: func(count int64) {
				if count < 9 || count > 11 {
					t.Fatalf("expected %d-%d, actual %d", 9, 11, count)
				}
			},
			d:    time.Millisecond * 100,
			done: createTimer(time.Millisecond * 1000),
		}},
		{"single incomplete iteration", args{
			f: func() {
				time.Sleep(time.Millisecond * 10)
			},
			callback: func(count int64) {
				if count < 90 || count > 105 {
					t.Fatalf("expected %d-%d, actual %d", 90, 105, count)
				}
			},
			d:    time.Millisecond * 5000,
			done: createTimer(time.Millisecond * 1000),
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Benchmark(tt.args.f, tt.args.callback, tt.args.d, tt.args.done)
		})
	}
}
