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
