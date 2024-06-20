// The Open/Closed Principle (OCP)  states that software entities ( structure, functions, etc.) should be open for extension but closed for modification.
// This means you should be able to add new functionality to a system by adding new code rather than modifying existing code.
// open for extension but closed for modification 

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

/*
Open for Extension
Your Shape interface allows you to add new shapes (e.g., Rectangle, Circle) without modifying the existing code. You can define new types that implement the Shape interface and provide their own implementation of the Area method.

For example, if you wanted to add a Triangle shape, you could do so without changing any of the existing Rectangle or Circle code

Closed for Modification
Your existing code, including the PrintArea function, does not need to be changed to accommodate new shapes. The PrintArea function works with any type that implements the Shape interface, so it remains unchanged even when new shapes are added.

*/