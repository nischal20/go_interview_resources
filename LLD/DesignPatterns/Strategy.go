//This pattern is useful when you need to select an algorithm at runtime and want to make your code more flexible and maintainable by adhering to the open/closed principle (open for extension, closed for modification).
package main
import "fmt"


type Strategy interface {
  Execute(a, b int) int
}
type Context struct {
  strategy Strategy
}
func (c *Context) AddStratergy(s Strategy) {
  c.strategy = s
}
func (c *Context) ExecuteStratergy(a, b int) int {
  return c.strategy.Execute(a, b)
}


type Add struct {}


func (add *Add) Execute(a, b int) int {
  return a + b
}
type Subtract struct {}
func (s *Subtract) Execute(a, b int) int {
  return a - b
}
func main() {
  c := Context{}
  c.AddStratergy(&Add{})
  fmt.Println(c.strategy.Execute(3, 5))


  c.AddStratergy(&Subtract{})
  fmt.Println(c.strategy.Execute(3, 5))
}

/*
Your observation is correct: the AddStratergy method in the Context struct is defined to accept a Strategy interface, which doesn't explicitly require a pointer. However, in Go, interfaces are implemented implicitly, and both pointer and value types can implement an interface.

In your current implementation, the Execute method is defined with pointer receivers (*Add and *Subtract), which means the instances you pass must be pointers to satisfy the interface.

Here's a breakdown:

Interface Definition:

go
Copy code
type Strategy interface {
  Execute(a, b int) int
}
Structs Implementing the Interface with Pointer Receivers:

go
Copy code
type Add struct {}
func (add *Add) Execute(a, b int) int {
  return a + b
}

type Subtract struct {}
func (s *Subtract) Execute(a, b int) int {
  return a - b
}
Since Execute is defined with a pointer receiver, the interface implementation requires the method to be called on a pointer to the struct (*Add and *Subtract).

Context Struct:
go
Copy code
type Context struct {
  strategy Strategy
}

func (c *Context) AddStratergy(s Strategy) {
  c.strategy = s
}

func (c *Context) ExecuteStratergy(a, b int) int {
  return c.strategy.Execute(a, b)
}
Even though AddStratergy takes a Strategy interface, the actual implementation requires pointers because the methods are defined on pointer receivers.

Possible Adjustments
If you want to allow passing values directly without requiring pointers, you can change the method receivers to value receivers:

go
Copy code
type Add struct {}
func (add Add) Execute(a, b int) int {
  return a + b
}

type Subtract struct {}
func (s Subtract) Execute(a, b int) int {
  return a - b
}
With these changes, you can pass values directly:

go
Copy code
func main() {
  c := Context{}
  c.AddStratergy(Add{})
  fmt.Println(c.ExecuteStratergy(3, 5)) // Output: 8

  c.AddStratergy(Subtract{})
  fmt.Println(c.ExecuteStratergy(3, 5)) // Output: -2
}
In this version, Add{} and Subtract{} are values, and the methods operate on copies of these values.

Why Use Pointers Originally?
Performance: If the struct is large, passing pointers is more efficient than copying large structs.
Mutability: If methods need to modify the state of the struct, pointers are required.
Conclusion
The initial design with pointer receivers is a common pattern for efficiency and flexibility, but if the struct is small and immutability is not a concern, using value receivers can simplify the code. In either case, Go's interfaces are flexible enough to work with both approaches, but the methods' receiver type dictates whether you need to pass pointers or values.*/