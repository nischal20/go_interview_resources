package main

import "fmt"

type Bird interface {
    Fly()
}

type Sparrow struct{}

func (s Sparrow) Fly() {
    fmt.Println("Sparrow is flying")
}

type Penguin struct{}

func (p Penguin) Fly() {
    fmt.Println("Penguins can't fly")
}

func LetBirdFly(b Bird) {
    b.Fly()
}

func main() {
    sparrow := Sparrow{}
    penguin := Penguin{}

    LetBirdFly(sparrow)
    LetBirdFly(penguin) // This violates LSP since penguins can't fly
}
