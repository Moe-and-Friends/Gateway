package config

type Config struct {
	// AMQP (RabbitMQ) configuration
	Amqp struct {
		Username string `envconfig:"GATEWAY_AMQP_USERNAME" required:"true"`
		Password string `envconfig:"GATEWAY_AMQP_PASSWORD" required:"true"`
		Url      string `envconfig:"GATEWAY_AMQP_URL" requied:"true"`
		Port     string `envconfig:"GATEWAY_AMQP_PORT" default:"5672"`
		Vhost    string `envconfig:"GATEWAY_AMQP_VHOST" default:""`
		Exchange string `envconfig:"GATEWAY_AMQP_EXCHANGE" required:"true"`
	}

	// Redis (or Redis-compliant) configuration
	Redis struct {
		Url string `envconfig:"GATEWAY_REDIS_URL" required:"true"` // Require full URI so users can freely define schemes, etc.
	}
}
