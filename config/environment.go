package config

import (
	"os"
	"strings"
)

// Environment represents the current deployment environment
type Environment string

const (
	Development Environment = "development"
	Staging     Environment = "staging"
	Production  Environment = "production"
	Test        Environment = "test"
)

// GetEnvironment returns the current environment
func GetEnvironment() Environment {
	env := strings.ToLower(os.Getenv("ENVIRONMENT"))
	switch env {
	case "production":
		return Production
	case "staging":
		return Staging
	case "test":
		return Test
	default:
		return Development
	}
}

// IsProduction returns true if running in production
func IsProduction() bool {
	return GetEnvironment() == Production
}

// IsDevelopment returns true if running in development
func IsDevelopment() bool {
	return GetEnvironment() == Development
}

// IsStaging returns true if running in staging
func IsStaging() bool {
	return GetEnvironment() == Staging
}

// GetSwaggerHost returns the appropriate host for Swagger documentation
func GetSwaggerHost() string {
	env := GetEnvironment()
	switch env {
	case Production:
		return os.Getenv("PROD_HOST") // e.g., "api.yourdomain.com"
	case Staging:
		// If no staging host is set, fall back to production host
		if stagingHost := os.Getenv("STAGING_HOST"); stagingHost != "" {
			return stagingHost
		}
		return os.Getenv("PROD_HOST")
	default:
		return os.Getenv("DEV_HOST") // e.g., "localhost:8082"
	}
}

// GetSwaggerScheme returns the appropriate scheme for Swagger documentation
func GetSwaggerScheme() string {
	env := GetEnvironment()
	switch env {
	case Production, Staging:
		return "https"
	default:
		return "http"
	}
}
