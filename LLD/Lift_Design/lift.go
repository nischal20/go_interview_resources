package main

import (
	"fmt"
	"sync"
	"time"
)

// ElevatorStatus defines the state of the elevator
type ElevatorStatus string

const (
	MovingUp   ElevatorStatus = "MovingUp"
	MovingDown ElevatorStatus = "MovingDown"
	Idle       ElevatorStatus = "Idle"
)

// Pair holds the key-value pair
type Pair struct {
	key   int
	value int
}

// Elevator represents a single elevator
type Elevator struct {
	id           int
	currentFloor int
	status       ElevatorStatus
	targetFloors []int
	mu           sync.Mutex
}

// NewElevator creates a new elevator instance
func NewElevator(id int) *Elevator {
	return &Elevator{
		id:           id,
		currentFloor: 0,
		status:       Idle,
		targetFloors: []int{},
	}
}

// AddTargetFloor adds a new target floor to the elevator
func (e *Elevator) AddTargetFloor(floor int) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.targetFloors = append(e.targetFloors, floor)
	if e.status == Idle {
		go e.move()
	}
}

func (e *Elevator) move() {
	for {
		e.mu.Lock()
		if len(e.targetFloors) == 0 {
			e.status = Idle
			e.mu.Unlock()
			return
		}
		target := e.targetFloors[0]
		e.targetFloors = e.targetFloors[1:]
		e.mu.Unlock()

		if e.currentFloor < target {
			e.status = MovingUp
			for e.currentFloor < target {
				time.Sleep(1 * time.Second) // Simulate movement
				e.currentFloor++
				fmt.Printf("Elevator %d moving up to floor %d\n", e.id, e.currentFloor)
			}
		} else if e.currentFloor > target {
			e.status = MovingDown
			for e.currentFloor > target {
				time.Sleep(1 * time.Second) // Simulate movement
				e.currentFloor--
				fmt.Printf("Elevator %d moving down to floor %d\n", e.id, e.currentFloor)
			}
		}

		e.status = Idle
		fmt.Printf("Elevator %d reached floor %d\n", e.id, e.currentFloor)
	}
}

func (e *Elevator) Status() (int, ElevatorStatus) {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.currentFloor, e.status
}

// Floor represents a floor in the building
type Floor struct {
	number int
}

// NewFloor creates a new floor instance
func NewFloor(number int) *Floor {
	return &Floor{number: number}
}

func (f *Floor) Number() int {
	return f.number
}

// ElevatorController manages multiple elevators and floor requests
type ElevatorController struct {
	elevators []*Elevator
	mu        sync.Mutex
}

// NewElevatorController creates a new ElevatorController
func NewElevatorController(numElevators int) *ElevatorController {
	elevators := make([]*Elevator, numElevators)
	for i := 0; i < numElevators; i++ {
		elevators[i] = NewElevator(i)
	}
	return &ElevatorController{
		elevators: elevators,
	}
}

// RequestElevator handles floor requests and assigns an elevator
func (ec *ElevatorController) RequestElevator(floor int) {
	ec.mu.Lock()
	defer ec.mu.Unlock()

	// Find the best elevator to handle the request
	bestElevator := ec.elevators[0]
	for _, elevator := range ec.elevators {
		if elevator.status == Idle {
			bestElevator = elevator
			break
		}
	}

	fmt.Printf("Assigning elevator %d to floor %d\n", bestElevator.id, floor)
	bestElevator.AddTargetFloor(floor)
}

func (ec *ElevatorController) ElevatorStatus() {
	for _, elevator := range ec.elevators {
		currentFloor, status := elevator.Status()
		fmt.Printf("Elevator %d is at floor %d and is %s\n", elevator.id, currentFloor, status)
	}
}

func main() {
	controller := NewElevatorController(3)

	controller.RequestElevator(5)
	controller.RequestElevator(2)
	controller.RequestElevator(8)

	// Give some time for elevators to move
	for i := 0; i < 20; i++ {
		controller.ElevatorStatus()
	}
}
