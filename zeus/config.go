package zeus

import "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

const (
	DefaultGRPCPort   = "9090"
	DefaultRESTPort   = "8080"
	DefaultEnableCORS = true
	DefaultOnlyJSON   = true
)

type Config struct {
	gRPCPort         string
	restPort         string
	enableCORS       bool
	onlyJSON         bool
	restServeMuxOpts []runtime.ServeMuxOption
}

type ConfigFunc func(c *Config)

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

func AddRestServerMuxOpt(opt ...runtime.ServeMuxOption) ConfigFunc {
	return func(c *Config) {
		c.restServeMuxOpts = append(c.restServeMuxOpts, opt...)
	}
}

func generateConfig(args ...ConfigFunc) *Config {
	c := &Config{
		gRPCPort:   DefaultGRPCPort,
		restPort:   DefaultRESTPort,
		enableCORS: DefaultEnableCORS,
		onlyJSON:   DefaultOnlyJSON,
	}
	for i := range args {
		args[i](c)
	}
	return c
}
