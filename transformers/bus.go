package transformers

import (
	"log"
	"server/types"
	"time"
)

var secondsToKeepAlive = time.Second * 30

// Bus transforms the siri xml response from DFT to json.
func Bus(siri types.Siri) []types.Bus {
	buses := siri.ServiceDelivery.VehicleMonitoringDelivery.VehicleActivity

	cutOfTime := time.Now().Add(-secondsToKeepAlive)

	jsonBus := make([]types.Bus, len(buses))

	var counter uint8

	for _, bus := range buses {
		lastUpdated, err := time.Parse(time.RFC3339, bus.RecorderAtTime)

		if err != nil {
			log.Printf("Failed to parse bus last update time, error: %s", err.Error())
		} else {
			if lastUpdated.After(cutOfTime) {
				var newBus types.Bus

				newBus.ID = bus.MonitoredVehicleJourney.VehicleRef
				newBus.Route = types.BusRoute{
					ID:   bus.MonitoredVehicleJourney.LineRef,
					Name: bus.MonitoredVehicleJourney.PublishedLineName,
				}
				newBus.LastUpdated = lastUpdated
				newBus.Location = types.Coordinate{
					Longitude: bus.MonitoredVehicleJourney.VehicleLocation.Longitude,
					Latitude:  bus.MonitoredVehicleJourney.VehicleLocation.Latitude,
				}
				newBus.Bearing = bus.MonitoredVehicleJourney.Bearing

				jsonBus[counter] = newBus

				counter++
			}
		}
	}

	log.Printf(
		"%d/%d buses received within %d seconds",
		counter,
		len(buses),
		secondsToKeepAlive/time.Second,
	)

	return jsonBus[:counter]
}
