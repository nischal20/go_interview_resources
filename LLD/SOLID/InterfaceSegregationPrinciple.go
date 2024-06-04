package main

import "fmt"

type Printer interface {
    Print()
}

type Scanner interface {
    Scan()
}

type MultiFunctionDevice interface {
    Printer
    Scanner
}

type MultiFunctionPrinter struct{}

func (mfp MultiFunctionPrinter) Print() {
    fmt.Println("Printing...")
}

func (mfp MultiFunctionPrinter) Scan() {
    fmt.Println("Scanning...")
}

func main() {
    mfp := MultiFunctionPrinter{}
    mfp.Print()
    mfp.Scan()
}
