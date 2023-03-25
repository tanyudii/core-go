package auth

type Config struct{}

type ConfigFunc func(c *Config)

func generateConfig(args ...ConfigFunc) *Config {
	c := &Config{}
	for i := range args {
		args[i](c)
	}
	return c
}
