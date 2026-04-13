package main

import "fmt"

/*
PRODUCT
*/
type House struct {
	Doors   int
	Windows int
	Garage  bool
}

/*
BUILDER
*/
type HouseBuilder struct {
	doors   int
	windows int
	garage  bool
}

func NewHouseBuilder() *HouseBuilder {
	return &HouseBuilder{}
}

func (b *HouseBuilder) SetDoors(doors int) *HouseBuilder {
	b.doors = doors
	return b
}

func (b *HouseBuilder) SetWindows(windows int) *HouseBuilder {
	b.windows = windows
	return b
}

func (b *HouseBuilder) SetGarage(garage bool) *HouseBuilder {
	b.garage = garage
	return b
}

func (b *HouseBuilder) Build() *House {
	return &House{
		Doors:   b.doors,
		Windows: b.windows,
		Garage:  b.garage,
	}
}

/*
CLIENT
*/
func main() {
	house := NewHouseBuilder().
		SetDoors(2).
		SetWindows(4).
		SetGarage(true).
		Build()

	fmt.Printf("%+v\n", house)
}
