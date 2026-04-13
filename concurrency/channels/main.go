package main

import (
	"fmt"
	"strconv"
	"sync"
)

func printerWorker(wg *sync.WaitGroup, ch chan<- string, name string) {
	defer wg.Done()

	for i := 1; i <= 3; i++ {
		ch <- name + " : " + strconv.Itoa(i)
	}
}

func main() {
	ch := make(chan string)

	var wg sync.WaitGroup
	wg.Add(2)

	go printerWorker(&wg, ch, "printer worker one")
	go printerWorker(&wg, ch, "printer worker two")

	// close after all printer workers finish
	go func() {
		fmt.Println("waiting...")
		wg.Wait()
		close(ch)
	}()

	for val := range ch {
		fmt.Println("received:", val)
	}

	fmt.Println("done")
}
