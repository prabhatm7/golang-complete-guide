package main

import "fmt"

/*
PRODUCT
*/
type Vehicle interface {
	Type() string
}

/*
CONCRETE PRODUCTS
*/
type Car struct{}

func (c *Car) Type() string {
	return "Car"
}

type Bike struct{}

func (b *Bike) Type() string {
	return "Bike"
}

/*
FACTORY METHOD
*/
func NewVehicle(vehicleType string) Vehicle {
	switch vehicleType {
	case "car":
		return &Car{}
	case "bike":
		return &Bike{}
	default:
		return nil
	}
}

/*
CLIENT
*/
func main() {
	v1 := NewVehicle("car")
	v2 := NewVehicle("bike")

	fmt.Println(v1.Type())
	fmt.Println(v2.Type())
}
