package main

import (
	"fmt"
	"time"
)

func main() {
	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()

	timeout := time.After(1200 * time.Millisecond)
	count := 0

	for {
		select {
		case <-ticker.C:
			count++
			fmt.Println("tick", count)
		case <-timeout:
			fmt.Println("timeout done")
			return
		}
	}
}
