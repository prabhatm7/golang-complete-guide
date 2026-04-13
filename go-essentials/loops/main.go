package main

import "fmt"

func main() {

	// 1 Classic for loop
	fmt.Println("1. Classic for loop")
	for i := 1; i <= 5; i++ {
		fmt.Println(i)
	}

	// 2 for loop as while loop
	fmt.Println("\n2. for loop as while")
	j := 1
	for j <= 5 {
		fmt.Println(j)
		j++
	}

	// 3 Infinite loop
	fmt.Println("\n3. Infinite loop (break after condition)")
	k := 1
	for {
		fmt.Println(k)
		if k == 3 {
			break // exits infinite loop
		}
		k++
	}

	// 4 Range loop (slice)
	fmt.Println("\n4. Range loop over slice")
	nums := []int{10, 20, 30}
	for index, value := range nums {
		fmt.Printf("Index: %d, Value: %d\n", index, value)
	}

	// 5 Range loop (map)
	fmt.Println("\n5. Range loop over map")
	m := map[string]int{
		"apple":  2,
		"banana": 5,
	}

	for key, value := range m {
		fmt.Printf("%s -> %d\n", key, value)
	}

	//  6 Range loop (string - runes)
	fmt.Println("\n6. Range loop over string")
	str := "GoLang"
	for index, char := range str {
		fmt.Printf("Index: %d, Char: %c\n", index, char)
	}

	// 7 Loop with continue
	fmt.Println("\n7. Loop with continue (skip even numbers)")
	for i := 1; i <= 5; i++ {
		if i%2 == 0 {
			continue
		}
		fmt.Println(i)
	}

	// 8 Nested loop
	fmt.Println("\n8. Nested loop")
	for i := 1; i <= 3; i++ {
		for j := 1; j <= 3; j++ {
			fmt.Printf("i=%d j=%d\n", i, j)
		}
	}

	// 9 Labeled break (advanced but important)
	fmt.Println("\n9. Labeled break")
outer:
	for i := 1; i <= 3; i++ {
		for j := 1; j <= 3; j++ {
			if i == 2 && j == 2 {
				break outer
			}
			fmt.Printf("i=%d j=%d\n", i, j)
		}
	}
}
