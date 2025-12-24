package main

import (
	"fmt"
	"sync"
	"time"
)

// Task represents the work to be done
type Task struct {
	ID int
}

func main() {
	const numJobs = 10   // Total tasks to complete
	const numWorkers = 3 // Size of the worker pool

	jobs := make(chan Task, numJobs)
	results := make(chan int, numJobs)
	var wg sync.WaitGroup

	// 1. Start the Workers
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	// 2. Send Jobs to the channel
	for j := 1; j <= numJobs; j++ {
		jobs <- Task{ID: j}
	}
	close(jobs) // Closing tells workers no more jobs are coming

	// 3. Wait for all workers to finish in a separate goroutine
	go func() {
		wg.Wait()
		close(results) // Close results once all workers are done
	}()

	// 4. Collect results
	for res := range results {
		fmt.Printf("Result received: %d\n", res)
	}

	fmt.Println("All workers finished processing.")
}

// worker function: each worker runs this concurrently
func worker(id int, jobs <-chan Task, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		fmt.Printf("Worker %d started job %d\n", id, job.ID)

		// Simulate processing time
		time.Sleep(time.Second)

		fmt.Printf("Worker %d finished job %d\n", id, job.ID)
		results <- job.ID * 2 // Send result back
	}
}
