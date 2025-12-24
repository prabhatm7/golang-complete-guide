package main

import "fmt"

func main() {
	// string cannot be single quote
	fmt.Print("Hello World!")
}

// IMPORTANT:
// 1. package main is important for the go build command
// 2. go build command finds the module(go project) and then find the go main package
