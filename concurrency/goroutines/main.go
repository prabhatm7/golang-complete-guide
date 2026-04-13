package main

import (
	"fmt"
	"sync"
)

func printNumbers(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 1; i <= 3; i++ {
		fmt.Printf("goroutine %d -> %d\n", id, i)
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(3)

	go printNumbers(1, &wg)
	go printNumbers(2, &wg)
	go printNumbers(3, &wg)

	wg.Wait()
	fmt.Println("done")
}
