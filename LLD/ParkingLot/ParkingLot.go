package main

import (
	"fmt"
	"time"
)

// Vehicle interface
type Vehicle interface {
	GetVehicleType() string
}

type Car struct{}
type Bike struct{}
type Truck struct{}

func (c Car) GetVehicleType() string   { return "Car" }
func (b Bike) GetVehicleType() string  { return "Bike" }
func (t Truck) GetVehicleType() string { return "Truck" }

// ParkingSpot struct
type ParkingSpot struct {
	spotID   int
	spotType string
	occupied bool
}

func (ps *ParkingSpot) CanPark(vehicle Vehicle) bool {
	if ps.occupied {
		return false
	}
	switch ps.spotType {
	case "Truck":
		return true
	case "Car":
		return vehicle.GetVehicleType() == "Car" || vehicle.GetVehicleType() == "Bike"
	case "Bike":
		return vehicle.GetVehicleType() == "Bike"
	}
	return false
}

func (ps *ParkingSpot) ParkVehicle() {
	ps.occupied = true
}

func (ps *ParkingSpot) RemoveVehicle() {
	ps.occupied = false
}

// Ticket struct
type Ticket struct {
	vehicle   Vehicle
	spot      *ParkingSpot
	level     *ParkingLevel
	startTime time.Time
	endTime   time.Time
	fee       float64
}

func NewTicket(vehicle Vehicle, spot *ParkingSpot, level *ParkingLevel) *Ticket {
	return &Ticket{vehicle: vehicle, spot: spot, level: level, startTime: time.Now()}
}

func (t *Ticket) EndParking() {
	t.endTime = time.Now()
}

func (t *Ticket) CalculateFee() {
	duration := t.endTime.Sub(t.startTime).Hours()
	switch t.spot.spotType {
	case "Car":
		t.fee = duration * 10
	case "Bike":
		t.fee = duration * 5
	case "Truck":
		t.fee = duration * 20
	}
}

// ParkingLevel struct
type ParkingLevel struct {
	levelID  int
	carSpots   []*ParkingSpot
	bikeSpots  []*ParkingSpot
	truckSpots []*ParkingSpot
}

func NewParkingLevel(levelID int, carCapacity int, bikeCapacity int, truckCapacity int) *ParkingLevel {
	level := &ParkingLevel{levelID: levelID}

	level.carSpots = make([]*ParkingSpot, carCapacity)
	for i := range level.carSpots {
		level.carSpots[i] = &ParkingSpot{spotID: i + 1, spotType: "Car"}
	}

	level.bikeSpots = make([]*ParkingSpot, bikeCapacity)
	for i := range level.bikeSpots {
		level.bikeSpots[i] = &ParkingSpot{spotID: i + 1, spotType: "Bike"}
	}

	level.truckSpots = make([]*ParkingSpot, truckCapacity)
	for i := range level.truckSpots {
		level.truckSpots[i] = &ParkingSpot{spotID: i + 1, spotType: "Truck"}
	}

	return level
}

func (pl *ParkingLevel) FindSpot(vehicle Vehicle) *ParkingSpot {
	spotTypes := [][]*ParkingSpot{pl.bikeSpots, pl.carSpots, pl.truckSpots}

	for _, spots := range spotTypes {
		for _, spot := range spots {
			if spot.CanPark(vehicle) {
				return spot
			}
		}
	}

	return nil
}

// ParkingLot struct
type ParkingLot struct {
	levels  []*ParkingLevel
	tickets []*Ticket
}

func NewParkingLot() *ParkingLot {
	return &ParkingLot{}
}

func (pl *ParkingLot) AddLevel(level *ParkingLevel) {
	pl.levels = append(pl.levels, level)
}

func (pl *ParkingLot) Park(vehicle Vehicle) *Ticket {
	for _, level := range pl.levels {
		spot := level.FindSpot(vehicle)
		if spot != nil {
			spot.ParkVehicle()
			ticket := NewTicket(vehicle, spot, level)
			pl.tickets = append(pl.tickets, ticket)
			return ticket
		}
	}
	return nil
}

func (pl *ParkingLot) Unpark(ticket *Ticket) {
	ticket.EndParking()
	ticket.CalculateFee()
	ticket.spot.RemoveVehicle()
	fmt.Printf("Vehicle Type: %s, Level: %d, Spot: %d, Fee: %.2f\n", ticket.vehicle.GetVehicleType(), ticket.level.levelID, ticket.spot.spotID, ticket.fee)
}

func main() {
	parkingLot := NewParkingLot()

	// Adding levels with different capacities for car, bike, and truck spots
	parkingLot.AddLevel(NewParkingLevel(1, 10, 5, 3))
	parkingLot.AddLevel(NewParkingLevel(2, 15, 7, 2))
	parkingLot.AddLevel(NewParkingLevel(3, 20, 10, 5))

	car := Car{}
	bike := Bike{}
	truck := Truck{}

	ticket1 := parkingLot.Park(car)
	time.Sleep(2 * time.Second)
	parkingLot.Unpark(ticket1)

	ticket2 := parkingLot.Park(bike)
	time.Sleep(1 * time.Second)
	parkingLot.Unpark(ticket2)

	ticket3 := parkingLot.Park(truck)
	time.Sleep(3 * time.Second)
	parkingLot.Unpark(ticket3)

	// Attempt to park another bike when bike spots are full
	for i := 0; i < 5; i++ {
		parkingLot.Park(Bike{})
	}
	ticket4 := parkingLot.Park(bike) // This should park in a car spot or truck spot if car spots are full
	time.Sleep(2 * time.Second)
	parkingLot.Unpark(ticket4)
}
