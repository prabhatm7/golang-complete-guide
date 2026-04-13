package main

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	ID int
}

func main() {
	const numJobs = 5
	const numWorkers = 2

	jobs := make(chan Task, numJobs)
	results := make(chan int, numJobs)
	var wg sync.WaitGroup

	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	for j := 1; j <= numJobs; j++ {
		jobs <- Task{ID: j}
	}
	close(jobs)

	// close results when workers finish
	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		fmt.Printf("result: %d\n", res)
	}

	fmt.Println("done")
}

func worker(id int, jobs <-chan Task, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		fmt.Printf("worker %d started job %d\n", id, job.ID)

		time.Sleep(300 * time.Millisecond)

		fmt.Printf("worker %d finished job %d\n", id, job.ID)
		results <- job.ID * 2
	}
}
