package main

import "fmt"

func printNumbers(done chan bool) {
	for i := 1; i <= 5; i++ {
		fmt.Println(i)
	}
	done <- true // signal completion
}

func main() {
	done := make(chan bool)

	go printNumbers(done)

	<-done // wait for goroutine to finish
}
