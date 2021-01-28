package types

import (
	"encoding/xml"
)

// Siri the dft standard
type Siri struct {
	XMLName         xml.Name        `xml:"Siri"`
	ServiceDelivery ServiceDelivery `xml:"ServiceDelivery"`
}

// ServiceDelivery the dft
type ServiceDelivery struct {
	XMLName                   xml.Name                  `xml:"ServiceDelivery"`
	ResponseTimestamp         string                    `xml:"ResponseTimestamp"`
	ProducerRef               string                    `xml:"ProducerRef"`
	VehicleMonitoringDelivery VehicleMonitoringDelivery `xml:"VehicleMonitoringDelivery"`
}

// VehicleMonitoringDelivery the dft
type VehicleMonitoringDelivery struct {
	XMLName               xml.Name          `xml:"VehicleMonitoringDelivery"`
	ResponseTimestamp     string            `xml:"ResponseTimestamp"`
	RequestMessageRef     string            `xml:"RequestMessageRef"`
	ValidUntil            string            `xml:"ValidUntil"`
	ShortestPossibleCycle string            `xml:"ShortestPossibleCycle"`
	VehicleActivity       []VehicleActivity `xml:"VehicleActivity"`
}

// VehicleActivity the dft
type VehicleActivity struct {
	XMLName                 xml.Name                `xml:"VehicleActivity"`
	RecorderAtTime          string                  `xml:"RecordedAtTime"`
	ItemIdentifier          string                  `xml:"ItemIdentifier"`
	ValidUntilTime          string                  `xml:"ValidUntilTime"`
	MonitoredVehicleJourney MonitoredVehicleJourney `xml:"MonitoredVehicleJourney"`
}

// MonitoredVehicleJourney for dft
type MonitoredVehicleJourney struct {
	XMLName                  xml.Name        `xml:"MonitoredVehicleJourney"`
	LineRef                  string          `xml:"LineRef"`
	DirectionRef             string          `xml:"DirectionRef"`
	PublishedLineName        string          `xml:"PublishedLineName"`
	OperatorRef              string          `xml:"OperatorRef"`
	DestinationRef           string          `xml:"DestinationRef"`
	DestinationName          string          `xml:"DestinationName"`
	OriginAimedDepartureTime string          `xml:"OriginAimedDepartureTime"`
	VehicleLocation          VehicleLocation `xml:"VehicleLocation"`
	BlockRef                 string          `xml:"BlockRef"`
	Bearing                  float32         `xml:"Bearing"`
	VehicleJourneyRef        string          `xml:"VehicleJourneyRef"`
	VehicleRef               string          `xml:"VehicleRef"`
}

// VehicleLocation for dft
type VehicleLocation struct {
	XMLName   xml.Name `xml:"VehicleLocation"`
	Longitude float32  `xml:"Longitude"`
	Latitude  float32  `xml:"Latitude"`
}
