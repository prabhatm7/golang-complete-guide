package main

import (
	"fmt"
	"math"
)

// inflationRate
const inflation = 2.5

// calculateFutureValues calculates future value and real (inflation-adjusted) value
func calculateFutureValues(
	investmentAmount float64,
	expectedReturnRate float64,
	years float64,
) (float64, float64) {

	futureValue := investmentAmount * math.Pow(1+expectedReturnRate/100, years)
	futureRealValue := futureValue / math.Pow(1+inflation/100, years)

	return futureValue, futureRealValue
}

func main() {
	fmt.Println("Basic Calculator!")

	investmentAmount := 1000.0
	expectedReturnRate := 5.5
	years := 10.0

	fmt.Print("Enter years: ")
	_, err := fmt.Scan(&years)
	if err != nil {
		return
	}

	futureValue, futureRealValue :=
		calculateFutureValues(investmentAmount, expectedReturnRate, years)

	fmt.Printf("Expected Future Value : %.2f\n", futureValue)
	fmt.Printf("Expected Real Value   : %.4f\n", futureRealValue)
}
