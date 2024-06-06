//The Observer Pattern is a behavioural design pattern where an object (the subject) maintains a list of its dependents (observers) and notifies them of any state changes, usually by calling one of their methods. This pattern is useful for implementing distributed event-handling systems, where multiple objects need to stay updated with changes in another object.

package main

import (
	"fmt"
)

type Observer interface {
	Update(data string)
}
type Subject interface {
	Register(observer Observer)
	Deregister(observer Observer)
	NotifyAll()
}

type ConcreteSubject struct {
	observers []Observer
	data      string
}

func (s *ConcreteSubject) Register(observer Observer) {
	s.observers = append(s.observers, observer)
}

// Deregister method to remove an observer
func (s *ConcreteSubject) Deregister(observer Observer) {
	var index int
	for i, obs := range s.observers {
		if obs == observer {
			index = i
			break
		}
	}
	s.observers = append(s.observers[:index], s.observers[index+1:]...)
}

// NotifyAll method to notify all observers
func (s *ConcreteSubject) NotifyAll() {
	for _, observer := range s.observers {
		observer.Update(s.data)
	}
}

// SetData method to update the data and notify observers
func (s *ConcreteSubject) SetData(data string) {
	s.data = data
	s.NotifyAll()
}

// ConcreteObserver struct
type ConcreteObserver struct {
	id int
}

// Update method to update the observer with new data
func (c *ConcreteObserver) Update(data string) {
	fmt.Printf("Observer %d received data: %s\n", c.id, data)
}

func main() {
	// Create a subject
	subject := &ConcreteSubject{}

	// Create observers
	observer1 := &ConcreteObserver{id: 1}
	observer2 := &ConcreteObserver{id: 2}

	// Register observers with the subject
	subject.Register(observer1)
	subject.Register(observer2)

	// Change the data in the subject and notify observers
	subject.SetData("Hello, Observers!")

	// Deregister an observer and update data again
	subject.Deregister(observer1)
	subject.SetData("Observer 1 deregistered.")
}
