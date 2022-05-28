package configs

type Config struct {
	Host      string
	CacheSize int
}

var c *Config //nolint:gochecknoglobals

func GetConfig() Config {
	if c == nil {
		c = &Config{
			Host:      ":8080",
			CacheSize: 10 * 1024 * 1024, //nolint:gomnd
		}
	}

	return *c
}
