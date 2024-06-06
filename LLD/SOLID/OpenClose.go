// The Open/Closed Principle (OCP)  states that software entities ( structure, functions, etc.) should be open for extension but closed for modification.
// This means you should be able to add new functionality to a system by adding new code rather than modifying existing code.

package main

import "fmt"

type Shape interface {
    Area() float64
}

type Rectangle struct {
    Width, Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return 3.14 * c.Radius * c.Radius
}

func PrintArea(s Shape) {
    fmt.Println("Area:", s.Area())
}

func main() {
    r := Rectangle{Width: 5, Height: 10}
    c := Circle{Radius: 7}

    PrintArea(r)
    PrintArea(c)
}
