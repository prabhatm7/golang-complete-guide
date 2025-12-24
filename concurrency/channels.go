package main

import (
	"context"
	"fmt"
	"strconv"
	"sync"
)

func producer(ctx context.Context, wg *sync.WaitGroup, ch chan<- string, name string) {
	defer wg.Done()

	for i := 1; i <= 5; i++ {
		select {
		case <-ctx.Done():
			fmt.Println("producer done")
			return
		case ch <- name + " : " + strconv.Itoa(i):
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan string)

	var wg sync.WaitGroup
	wg.Add(2)

	go producer(ctx, &wg, ch, "producer one")
	go producer(ctx, &wg, ch, "producer two")

	// Close channel once ALL producers finish
	go func() {
		wg.Wait()
		close(ch)
	}()

	// Consumer (for–select event loop)
	for {
		select {
		case val, ok := <-ch:
			if !ok {
				fmt.Println("All producers finished. Exiting cleanly.")
				return
			}
			fmt.Println("Received:", val)
		case <-ctx.Done():
			fmt.Println("Context cancelled. Shutting down.")
			return
		}
	}
}
