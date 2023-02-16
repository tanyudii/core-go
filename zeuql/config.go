package zeuql

const (
	DefaultGraphQLPort      = "8080"
	DefaultEnableCORS       = true
	DefaultEnablePrometheus = true
	DefaultPrometheusPort   = "9800"
)

type Config struct {
	graphQLPort      string
	enableCORS       bool
	enablePrometheus bool
	prometheusPort   string
}

type ConfigFunc func(c *Config)

func generateConfig(args ...ConfigFunc) *Config {
	c := &Config{
		graphQLPort:      DefaultGraphQLPort,
		enableCORS:       DefaultEnableCORS,
		enablePrometheus: DefaultEnablePrometheus,
		prometheusPort:   DefaultPrometheusPort,
	}
	for i := range args {
		args[i](c)
	}
	return c
}

func GraphQLPort(p string) ConfigFunc {
	return func(c *Config) {
		c.graphQLPort = p
	}
}

func EnableCORS(cors bool) ConfigFunc {
	return func(c *Config) {
		c.enableCORS = cors
	}
}

func EnablePrometheus(p bool) ConfigFunc {
	return func(c *Config) {
		c.enablePrometheus = p
	}
}

func PrometheusPort(p string) ConfigFunc {
	return func(c *Config) {
		c.prometheusPort = p
	}
}
