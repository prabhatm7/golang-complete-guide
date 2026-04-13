package main

import (
	"fmt"
	"sync"
)

func generate(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n * n
		}
	}()
	return out
}

func fanOut(in <-chan int, workers int) []<-chan int {
	outs := make([]<-chan int, 0, workers)
	for i := 1; i <= workers; i++ {
		out := make(chan int)
		go func(id int, out chan<- int) {
			defer close(out)
			for n := range in {
				fmt.Printf("worker %d got %d\n", id, n)
				out <- n * n
			}
		}(i, out)
		outs = append(outs, out)
	}
	return outs
}

func fanIn(channels ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	wg.Add(len(channels))
	for _, ch := range channels {
		go func(c <-chan int) {
			defer wg.Done()
			for v := range c {
				out <- v
			}
		}(ch)
	}

	// close after all inputs finish
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	fmt.Println("== simple pipeline ==")
	stage1 := generate(1, 2, 3, 4)
	stage2 := square(stage1)
	for v := range stage2 {
		fmt.Println("pipeline result:", v)
	}

	fmt.Println("\n== fan-out + fan-in ==")
	input := generate(1, 2, 3, 4, 5, 6)
	workers := fanOut(input, 3)
	merged := fanIn(workers...)

	for v := range merged {
		fmt.Println("merged result:", v)
	}
}
