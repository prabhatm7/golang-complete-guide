package main

import "fmt"

type Coffee interface {
	Cost() int
	Description() string
}
type SimpleCoffee struct{}

func (sc *SimpleCoffee) Cost() int {
	return 5
}

func (sc *SimpleCoffee) Description() string {
	return "Simple Coffee Description"
}

type CoffeeDecorator struct {
	coffee Coffee
}

type MilkDecorator struct {
	coffeeDecorator CoffeeDecorator
}

func NewMilkDecorator(c Coffee) MilkDecorator {
	return MilkDecorator{
		coffeeDecorator: CoffeeDecorator{coffee: c},
	}
}

func (m *MilkDecorator) Description() string {
	return m.coffeeDecorator.coffee.Description() + ".. Milk added"
}

func (m *MilkDecorator) Cost() int {
	return m.coffeeDecorator.coffee.Cost() + 5
}
func main() {
	coffee := &SimpleCoffee{}

	coffeeDecorator := NewMilkDecorator(coffee)
	fmt.Println(coffeeDecorator.Cost())
	fmt.Println(coffeeDecorator.Description())
}
