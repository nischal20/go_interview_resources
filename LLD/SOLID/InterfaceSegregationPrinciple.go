//The Interface Segregation Principle (ISP) states that clients should not be forced to depend on interfaces they do not use. Instead of one large, general-purpose interface, multiple smaller, more specific interfaces should be created. This allows clients to depend only on the methods that are relevant to them.
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
