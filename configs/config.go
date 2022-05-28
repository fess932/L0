package configs

type Config struct {
	Host      string
	CacheSize int
	NatsHost  string
	Topic     string
	DB        string
}

var c *Config //nolint:gochecknoglobals

func GetConfig() Config {
	if c == nil {
		c = &Config{
			Host:      ":8080",
			CacheSize: 10 * 1024 * 1024, //nolint:gomnd
			NatsHost:  "nats://localhost:4222",
			Topic:     "order",
			DB:        "postgres://wb:wb@localhost:5432/wb?sslmode=disable",
		}
	}

	return *c
}
