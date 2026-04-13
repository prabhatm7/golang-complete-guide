package main

import (
	"fmt"
	"sync"
)

type Order struct {
	ID     int
	Status string
}

func main() {
	var wg sync.WaitGroup

	orderChan := make(chan *Order)
	processedChan := make(chan *Order)

	wg.Add(3)

	go func() {
		defer wg.Done()
		defer close(orderChan)
		for _, order := range generateOrders(5) {
			orderChan <- order
		}
		fmt.Println("orders sent")
	}()

	go processOrders(&wg, orderChan, processedChan)

	go func() {
		defer wg.Done()
		for order := range processedChan {
			fmt.Printf("order %d -> %s\n", order.ID, order.Status)
		}
		fmt.Println("results done")
	}()

	wg.Wait()
	fmt.Println("all operations completed")
}

func generateOrders(count int) []*Order {
	orders := make([]*Order, count)
	for i := range orders {
		orders[i] = &Order{ID: i + 1, Status: "pending"}
	}
	return orders
}

func processOrders(wg *sync.WaitGroup, inChan <-chan *Order, outChan chan<- *Order) {
	defer wg.Done()
	defer close(outChan)

	for order := range inChan {
		order.Status = "processed"
		outChan <- order
	}
}
