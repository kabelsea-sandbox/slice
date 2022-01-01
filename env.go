package slice

import (
	"os"
)

// Env is a environment of the application.
type Env string

const (
	// Dev environment.
	EnvDev Env = "dev"
	// Test environment.
	EnvTest Env = "test"
	// Prod environment. Default.
	EnvProd Env = "prod"
)

// IsDevelopment returns true if project in development environment.
func (e Env) IsDevelopment() bool {
	return e == EnvDev
}

// IsProduction returns true if project in production environment.
func (e Env) IsProduction() bool {
	return e == EnvProd
}

// IsTesting returns true if project in testing environment.
func (e Env) IsTesting() bool {
	return e == EnvTest
}

var (
	// getEnv used in testing purposes
	getEnv = os.Getenv
)
