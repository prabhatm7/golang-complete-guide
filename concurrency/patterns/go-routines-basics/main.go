package main

import (
	"fmt"
	"sync"
)

func task(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("task %d running\n", id)
}

func main() {
	var wg sync.WaitGroup
	wg.Add(3)

	go task(1, &wg)
	go task(2, &wg)
	go task(3, &wg)

	wg.Wait()
	fmt.Println("main exited")

}
