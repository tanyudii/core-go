package zeus

import "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

const (
	DefaultGRPCPort           = "9090"
	DefaultRESTPort           = "8080"
	DefaultEnableCORS         = true
	DefaultOnlyJSON           = true
	DefaultEnablePrometheus   = true
	DefaultPrometheusPort     = "9800"
	DefaultRegisterReflection = true
)

type Config struct {
	gRPCPort           string
	restPort           string
	enableCORS         bool
	onlyJSON           bool
	enablePrometheus   bool
	prometheusPort     string
	registerReflection bool
	restServeMuxOpts   []runtime.ServeMuxOption
}

type ConfigFunc func(c *Config)

func generateConfig(args ...ConfigFunc) *Config {
	c := &Config{
		gRPCPort:           DefaultGRPCPort,
		restPort:           DefaultRESTPort,
		enableCORS:         DefaultEnableCORS,
		onlyJSON:           DefaultOnlyJSON,
		enablePrometheus:   DefaultEnablePrometheus,
		prometheusPort:     DefaultPrometheusPort,
		registerReflection: DefaultRegisterReflection,
	}
	for i := range args {
		args[i](c)
	}
	return c
}

func GRPCPort(p string) ConfigFunc {
	return func(c *Config) {
		c.gRPCPort = p
	}
}

func RESTPort(p string) ConfigFunc {
	return func(c *Config) {
		c.restPort = p
	}
}

func EnableCORS(cors bool) ConfigFunc {
	return func(c *Config) {
		c.enableCORS = cors
	}
}

func OnlyJSON(j bool) ConfigFunc {
	return func(c *Config) {
		c.onlyJSON = j
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

func RegisterReflection(r bool) ConfigFunc {
	return func(c *Config) {
		c.registerReflection = r
	}
}

func AddRestServerMuxOpt(opt ...runtime.ServeMuxOption) ConfigFunc {
	return func(c *Config) {
		c.restServeMuxOpts = append(c.restServeMuxOpts, opt...)
	}
}
