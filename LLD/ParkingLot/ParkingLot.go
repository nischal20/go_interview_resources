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
    spotID    int
    spotType  string
    occupied  bool
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
    vehicle    Vehicle
    spot       *ParkingSpot
    level      *ParkingLevel
    startTime  time.Time
    endTime    time.Time
    fee        float64
}

func NewTicket(vehicle Vehicle, spot *ParkingSpot, level *ParkingLevel) *Ticket {
    return &Ticket{vehicle: vehicle, spot: spot, level: level, startTime: time.Now()}
}

func (t *Ticket) EndParking() {
    t.endTime = time.Now()
}

func (t *Ticket) CalculateFee() {
    duration := t.endTime.Sub(t.startTime).Hours()
    switch t.vehicle.GetVehicleType() {
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
    levelID int
    spots   []*ParkingSpot
}

func NewParkingLevel(levelID int, capacity int, spotType string) *ParkingLevel {
    spots := make([]*ParkingSpot, capacity)
    for i := range spots {
        spots[i] = &ParkingSpot{spotID: i + 1, spotType: spotType}
    }
    return &ParkingLevel{levelID: levelID, spots: spots}
}

func (pl *ParkingLevel) FindSpot(vehicle Vehicle) *ParkingSpot {
    for _, spot := range pl.spots {
        if spot.CanPark(vehicle) {
            return spot
        }
    }
    return nil
}

// ParkingLot struct
type ParkingLot struct {
    levels []*ParkingLevel
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

    // Adding levels
    parkingLot.AddLevel(NewParkingLevel(1, 10, "Car"))
    parkingLot.AddLevel(NewParkingLevel(2, 5, "Truck"))
    parkingLot.AddLevel(NewParkingLevel(3, 15, "Bike"))

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
}
