package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Order struct {
	ID     int
	Status string
	mu     sync.Mutex // Mutex for individual order protection
}

func main() {
	var wg sync.WaitGroup
	orders := generateOrders(20)

	// Channels: Communication conduits between goroutines
	// Using a buffered channel with capacity 20
	orderChan := make(chan *Order, 20)
	processedChan := make(chan *Order, 20)

	wg.Add(3) // We will spawn 3 main concurrent processes

	// 1. Goroutine for Generating/Sending Orders
	go func() {
		defer wg.Done()
		defer close(orderChan) // Always close channel from sender side
		for _, order := range orders {
			orderChan <- order
		}
		fmt.Println("Done with generating orders [00:26:51]")
	}()

	// 2. Goroutine for Processing Orders
	go processOrders(&wg, orderChan, processedChan)

	// 3. Goroutine for Handling Results (using Select)
	go func() {
		defer wg.Done()
		for {
			select {
			case processedOrder, ok := <-processedChan:
				if !ok {
					fmt.Println("Processing Channel closed")
					return
				}
				fmt.Printf("Processed Order ID: %d, Status: %s\n", processedOrder.ID, processedOrder.Status)

			}
		}
	}()

	wg.Wait() // Blocks until the counter becomes zero [00:13:04]
	fmt.Println("All operations completed")
}

// generateOrders creates a slice of order pointers
func generateOrders(count int) []*Order {
	orders := make([]*Order, count)
	for i := 0; i < count; i++ {
		orders[i] = &Order{
			ID:     i + 1,
			Status: "Pending",
		}
	}
	return orders
}

// processOrders demonstrates receive-only and send-only channels [00:29:21]
func processOrders(wg *sync.WaitGroup, inChan <-chan *Order, outChan chan<- *Order) {
	defer wg.Done()
	defer close(outChan)

	for order := range inChan {
		// Simulate work
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)

		// Mutex usage to protect data [00:17:37]
		order.mu.Lock()
		order.Status = "Processed"
		order.mu.Unlock()

		outChan <- order
	}
}
