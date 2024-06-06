//Liskov Substitution Principle (LSP) states that objects of a superclass should be replaceable with objects of a subclass without affecting the correctness of the program. In other words, if S is a subtype of T, then objects of type T may be replaced with objects of type S without altering any of the desirable properties of the program 
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
