package config

type Config struct {
	DatabaseDSN string `envconfig:"DATABASE_DSN"`
	Addr        string `envconfig:"ADDR"`
}
