package common

import (
	"os"
	"strconv"
)

func GetPortFromEnvVar(defaultValue int) int {
	portEnv := os.Getenv("PORT")

	port := defaultValue
	if portEnv != "" {
		parsedPort, err := strconv.Atoi(portEnv)
		if err == nil {
			port = parsedPort
		}
	}

	return port
}
