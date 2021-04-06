package util

import (
	"fmt"
	"os"
)

// VerifyEnvSet verifies all environment variables are set and logs to stderr
// any which are not set
func VerifyEnvSet() bool {
	ok := true

	envVariables := []string{
		"DFT_SECRET", "MAPBOX_TOKEN", "DATABASE_URL", "ADMIN_TOKEN",
	}

	for _, envVar := range envVariables {
		if _, isPresent := os.LookupEnv(envVar); isPresent == false {
			fmt.Fprintln(os.Stderr, "No DFT_SECRET provided")
			ok = false
		}
	}

	return ok
}