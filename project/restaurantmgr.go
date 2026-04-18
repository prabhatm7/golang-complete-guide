package main

// Restaurant Manager
// Customer gives a Orders -> Kitchen prepares orders -> deliver the Orders to the customer -> finally generate the report
import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type OrderStatus string 

var mu sync.Mutex
var OrderCount, deliverCount = 0, 0

const (
	NewOrder OrderStatus = "NEW_ORDER"
	OrderPreparing OrderStatus = "ORDER_PREPARING"
	OrderPrepared OrderStatus = "ORDER_PREPARED"
	OrderDelivered OrderStatus = "ORDER_DELIVERED"
)

type Order struct {
	ID int
	CustomerName string
	OrderStatus OrderStatus
}


func getCustomerOrder(orderChan chan <- Order, wg *sync.WaitGroup) {
	// populate order channel every random time for total 50 orders
	defer wg.Done()
	defer close(orderChan)

	for i:= 1 ; i <= 10; i++ {
		// sleep at random time
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		fmt.Printf("Customer %s has placed order id : %d\n", "abcd"+ strconv.Itoa(i), i)
		orderChan <- Order{
			ID: i,
			CustomerName: "abcd"+ strconv.Itoa(i),
			OrderStatus: NewOrder,
		}
	}
}

func prepareOrder(orderChan <-chan Order, deliveryChan chan<- Order, chefNum int, wg *sync.WaitGroup) {
	defer wg.Done()

	for order := range orderChan {
		mu.Lock()
		OrderCount++
		mu.Unlock()

		order.OrderStatus = OrderPreparing

		deliveryChan <- order
		fmt.Printf("Chef Number %d preparing order id %d\n", chefNum, order.ID)
	}

}

func deliverOrder(deliveryChan <-chan Order, valetNum int, wg *sync.WaitGroup) {
	defer wg.Done()

	for order := range deliveryChan {
		mu.Lock()
		deliverCount++
		// order.OrderStatus = OrderDelivered
		mu.Unlock()
		fmt.Printf("Order Number %d delivered by valet id %d\n", order.ID, valetNum)
	}
}

func main() {
	var wg sync.WaitGroup
	var chefWg sync.WaitGroup
	maxOrderCap := 10
	chefCount := 3

	maxDeliveryCap := 5
	valetCount := 3

	orderChan := make(chan Order, maxOrderCap)
	deliveryChan := make(chan Order, maxDeliveryCap)

	wg.Add(1)
	go getCustomerOrder(orderChan, &wg)
	
	// chef workers
	for i:= 1; i <= chefCount; i++ {
		wg.Add(1)
		chefWg.Add(1)
		go func(chefNum int) {
			defer chefWg.Done()
			prepareOrder(orderChan, deliveryChan, chefNum, &wg)
		}(i)
	}

	go func() {
		chefWg.Wait()
		close(deliveryChan)
	}()

	// deliver orders
	for i:= 1; i <= valetCount; i++ {
		wg.Add(1)
		go deliverOrder(deliveryChan, i, &wg)
	}

	wg.Wait()

	fmt.Printf("Total Orders Prepared : %d\n", OrderCount)
	fmt.Printf("Total Orders Delivered: %d\n", deliverCount)
}