package types

import "time"

// Bus containing information about a bus
type Bus struct {
	ID          string
	Route       BusRoute
	Location    Coordinate
	Bearing     float32
	LastUpdated time.Time
}

// BusRoute contains information about the bus route
type BusRoute struct {
	ID   string
	Name string
}
