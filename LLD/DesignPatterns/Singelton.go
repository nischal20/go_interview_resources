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

// Changed method receiver from (s *Singleton) to () to make it a package-level function
func GetInstance() *Singleton {
	mu.Lock()
	defer mu.Unlock()
	if instance == nil {
		instance = &Singleton{data: "data"}
	}
	return instance
}

func main() {
	// Access the singleton instance directly through the package-level function
	fmt.Println(GetInstance().data)
	b := GetInstance()
	fmt.Println(b.data)

	// Modifying the singleton data to demonstrate it is the same instance
	b.data = "new data"
	fmt.Println(GetInstance().data)
}
