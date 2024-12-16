package errors

import (
	"errors"
	"fmt"
)

// Environment errors
var ErrNoEnvFile = errors.New("no .env file found, relying on system environment variables")

func EnvVariableNotSet(variableName string) error {
	return fmt.Errorf("environment variable %s is required but not set", variableName)
}
