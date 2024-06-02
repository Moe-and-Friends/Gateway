package main

type Config struct {
	// Redis (or Redis-compliant) configuration
	Redis struct {
		Url string `envconfig:"GATEWAY_REDIS_URL" required:"true"` // Require full URI so users can freely define schemes, etc.
	}
}
