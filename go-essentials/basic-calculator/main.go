package main

import (
	"fmt"
	"math"
)

const inflationRate = 2.5

func main() {
	fmt.Println("Basic Calculator!")

	var investmentAmount float64 = 1000
	var expectedReturnRate float64 = 5.5
	var years float64 = 10

	fmt.Print("Enter years: ")
	_, err := fmt.Scan(&years)
	if err != nil {
		return
	}

	var futureValue = investmentAmount * math.Pow(1+expectedReturnRate/100, years)
	var futureRealValue = futureValue / math.Pow(1+inflationRate/100, years)

	fmt.Printf("Expected Future Value : %f\n", futureValue)
	fmt.Printf("Expected Real : %.4f\n", futureRealValue)
}
