package zeuql

const (
	DefaultGraphQLPort      = "8080"
	DefaultGraphQLPath      = "/graphql"
	DefaultPlaygroundPath   = "/playground"
	DefaultEnablePlayground = true
	DefaultEnableCORS       = true
	DefaultEnablePrometheus = true
	DefaultPrometheusPort   = "9800"
)

type Config struct {
	graphQLPort      string
	graphQLPath      string
	playgroundPath   string
	enableCORS       bool
	enablePlayground bool
	enablePrometheus bool
	prometheusPort   string
}

type ConfigFunc func(c *Config)

func generateConfig(args ...ConfigFunc) *Config {
	c := &Config{
		graphQLPort:      DefaultGraphQLPort,
		graphQLPath:      DefaultGraphQLPath,
		playgroundPath:   DefaultPlaygroundPath,
		enablePlayground: DefaultEnablePlayground,
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

func GraphQLPath(p string) ConfigFunc {
	return func(c *Config) {
		c.graphQLPath = p
	}
}

func PlaygroundPath(p string) ConfigFunc {
	return func(c *Config) {
		c.playgroundPath = p
	}
}

func EnablePlayground(e bool) ConfigFunc {
	return func(c *Config) {
		c.enablePlayground = e
	}
}

func EnableCORS(e bool) ConfigFunc {
	return func(c *Config) {
		c.enableCORS = e
	}
}

func EnablePrometheus(e bool) ConfigFunc {
	return func(c *Config) {
		c.enablePrometheus = e
	}
}

func PrometheusPort(p string) ConfigFunc {
	return func(c *Config) {
		c.prometheusPort = p
	}
}
