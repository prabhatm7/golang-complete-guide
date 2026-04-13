package main

import "fmt"

func main() {
	var num int
	var choice int

	fmt.Print("Enter a number: ")
	_, err := fmt.Scan(&num)
	if err != nil {
		return
	}

	for {
		fmt.Println("\nChoose an option:")
		fmt.Println("1. Check even or odd")
		fmt.Println("2. Check positive, negative, or zero")
		fmt.Println("3. Print numbers from 1 to N (skip multiples of 3)")
		fmt.Println("4. Exit")

		fmt.Print("Enter choice: ")
		fmt.Scan(&choice)

		switch choice {
		case 1:
			if num%2 == 0 {
				fmt.Println("Number is EVEN")
			} else {
				fmt.Println("Number is ODD")
			}

		case 2:
			if num > 0 {
				fmt.Println("Number is POSITIVE")
			} else if num < 0 {
				fmt.Println("Number is NEGATIVE")
			} else {
				fmt.Println("Number is ZERO")
			}

		case 3:
			fmt.Println("Printing numbers:")
			for i := 1; i <= num; i++ {
				if i%3 == 0 {
					continue // skip multiples of 3
				}
				fmt.Print(i, " ")
			}
			fmt.Println()

		case 4:
			fmt.Println("Exiting program...")
			break // breaks switch
		default:
			fmt.Println("Invalid choice")
		}

		if choice == 4 {
			break // breaks for loop
		}
	}
}
