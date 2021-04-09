package controllers

import (
	"os"
	"server/models"
	"server/types"
	"server/transformers"
	"fmt"
)

// GetBusLocations get the bus locations within the a box of two coordinates
func GetBusLocations(
	topLeftCoordinate types.Coordinate,
	bottomRightCoordinate types.Coordinate,
) ([]types.Bus, error) {
	siriBuses, err := models.GetBusLocation(topLeftCoordinate, bottomRightCoordinate)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)

		return nil, err
	}

	return transformers.Bus(siriBuses), nil
}