package main

import (
	"fmt"
	"github.com/iskorotkov/chaos-io-stress/pkg/bench"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	out := make(chan os.Signal, 1)
	signal.Notify(out, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan struct{})
	go func() {
		<-out
		done <- struct{}{}
	}()

	filename := "dummy-file"
	writeDummyFile(filename, "content")
	defer func() {
		// Can't delete file immediately as it's not closed yet
		time.Sleep(time.Millisecond * 100)

		err := os.Remove(filename)
		if err != nil {
			panic(err)
		}
	}()

	bench.Benchmark(func() {
		content := readDummyFile(filename)
		writeDummyFile(filename, content)
	}, func(count int64) {
		fmt.Printf("executed %d times\n", count)
	}, time.Second, done)
}

func readDummyFile(path string) string {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("couldn't read file '%s'", path)
		panic(err)
	}

	return string(content)
}

func writeDummyFile(path string, content string) {
	err := ioutil.WriteFile(path, []byte(content), os.ModePerm)
	if err != nil {
		log.Printf("couldn't write file '%s'", path)
		panic(err)
	}
}
