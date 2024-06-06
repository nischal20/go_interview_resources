//In summary, the singleton pattern ensures that there is only one instance of the singleton struct. When you modify the data through one reference (s1), it is reflected in all other references (s2), since they all point to the same instance.
package main


import (
  "fmt"
  "sync"
)
type Singleton struct {
  data string
}


var instance *Singleton
var mu sync.Mutex
func GetInstance() *Singleton {
  mu.Lock()
  defer mu.Unlock()
  if instance == nil {
     instance = Singleton{data: "daata"}
  }
  return instance
}


func main() {
  a := GetInstance()
  fmt.Println(a.data)
  b := GetInstance()
  fmt.Println(b.data)
}
