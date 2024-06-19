package main

import (
    "fmt"
)

// ElevatorState interface
type ElevatorState interface {
    Move(e *Elevator)
    Stop(e *Elevator)
    Request(e *Elevator, floor int)
}

// Concrete States
type MovingUp struct{}

func (m *MovingUp) Move(e *Elevator) {
    e.currentFloor++
    fmt.Println("Elevator moving up to floor", e.currentFloor)
}

func (m *MovingUp) Stop(e *Elevator) {
    fmt.Println("Elevator stopping at floor", e.currentFloor)
    e.state = &Idle{}
}

func (m *MovingUp) Request(e *Elevator, floor int) {
    fmt.Println("Request received to move to floor", floor)
}

type MovingDown struct{}

func (m *MovingDown) Move(e *Elevator) {
    e.currentFloor--
    fmt.Println("Elevator moving down to floor", e.currentFloor)
}

func (m *MovingDown) Stop(e *Elevator) {
    fmt.Println("Elevator stopping at floor", e.currentFloor)
    e.state = &Idle{}
}

func (m *MovingDown) Request(e *Elevator, floor int) {
    fmt.Println("Request received to move to floor", floor)
}

type Idle struct{}

func (i *Idle) Move(e *Elevator) {
    // Implement Idle move behavior
}

func (i *Idle) Stop(e *Elevator) {
    // Implement Idle stop behavior
}

func (i *Idle) Request(e *Elevator, floor int) {
    if floor > e.currentFloor {
        e.state = &MovingUp{}
    } else {
        e.state = &MovingDown{}
    }
    e.Move()
}

// Elevator struct
type Elevator struct {
    currentFloor int
    state        ElevatorState
}

func (e *Elevator) Move() {
    e.state.Move(e)
}

func (e *Elevator) Stop() {
    e.state.Stop(e)
}

func (e *Elevator) Request(floor int) {
    e.state.Request(e, floor)
}

// ElevatorController struct
type ElevatorController struct {
    elevators []*Elevator
}

func (ec *ElevatorController) RequestElevator(floor int) {
    bestElevator := ec.FindBestElevator(floor)
    bestElevator.Request(floor)
}

func (ec *ElevatorController) FindBestElevator(floor int) *Elevator {
    // Implement scheduling strategy here
    // For simplicity, return the first available elevator
    for _, elevator := range ec.elevators {
        if _, ok := elevator.state.(*Idle); ok {
            return elevator
        }
    }
    return ec.elevators[0] // Default to the first elevator if none are idle
}

// Observer interface
type Observer interface {
    Notify(elevator *Elevator, event string)
}

type Display struct{}

func (d *Display) Notify(elevator *Elevator, event string) {
    fmt.Printf("Display: Elevator at floor %d, event: %s\n", elevator.currentFloor, event)
}

type ElevatorSystem struct {
    controller *ElevatorController
    observers  []Observer
}

func (es *ElevatorSystem) AddObserver(observer Observer) {
    es.observers = append(es.observers, observer)
}

func (es *ElevatorSystem) NotifyObservers(elevator *Elevator, event string) {
    for _, observer := range es.observers {
        observer.Notify(elevator, event)
    }
}

// Example usage
func main() {
    elevator1 := &Elevator{currentFloor: 0, state: &Idle{}}
    elevator2 := &Elevator{currentFloor: 0, state: &Idle{}}
    controller := &ElevatorController{elevators: []*Elevator{elevator1, elevator2}}
    system := &ElevatorSystem{controller: controller}

    display := &Display{}
    system.AddObserver(display)

    controller.RequestElevator(5)
}
