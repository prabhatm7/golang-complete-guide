package main

import (
	"fmt"
	"time"
)

func someFunc(num int) {
	fmt.Printf("id : %v\n", num)
}
func main() {
	go someFunc(1)
	go someFunc(2)
	go someFunc(3)
	time.Sleep(1 * time.Second)
	fmt.Printf("main exited!!")

}

// golang follows the fork join pattern
