package controllers

import (
	"os"
	"fmt"
)

type sqlDB interface {
	Ping() error
}

// HealthCheck tests the current health of the service
func HealthCheck(db sqlDB) bool {
	if err := db.Ping(); err != nil {
		fmt.Fprintln(os.Stderr, err)

		return false
	}

	return true
}