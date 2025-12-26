package main

import "fmt"

/*
STRATEGY
*/
type PaymentStrategy interface {
	Pay(amount int)
}

/*
CONCRETE STRATEGIES
*/
type CreditCard struct{}

func (c *CreditCard) Pay(amount int) {
	fmt.Println("Paid", amount, "using Credit Card")
}

type UPI struct{}

func (u *UPI) Pay(amount int) {
	fmt.Println("Paid", amount, "using UPI")
}

/*
CONTEXT
*/
type PaymentContext struct {
	strategy PaymentStrategy
}

func (p *PaymentContext) SetStrategy(s PaymentStrategy) {
	p.strategy = s
}

func (p *PaymentContext) Pay(amount int) {
	p.strategy.Pay(amount)
}

func main() {
	ctx := &PaymentContext{}

	ctx.SetStrategy(&CreditCard{})
	ctx.Pay(100)

	ctx.SetStrategy(&UPI{})
	ctx.Pay(200)
}
