package main

import "fmt"

/*
TARGET
*/
type Printer interface {
	Print()
}

/*
ADAPTEE
*/
type OldPrinter struct{}

func (o *OldPrinter) PrintOld() {
	fmt.Println("Printing using old printer")
}

/*
ADAPTER
*/
type PrinterAdapter struct {
	oldPrinter *OldPrinter
}

func (p *PrinterAdapter) Print() {
	p.oldPrinter.PrintOld()
}

func main() {
	old := &OldPrinter{}
	adapter := &PrinterAdapter{oldPrinter: old}

	adapter.Print()
}
