package auth

type MapPublicRoutes map[string]bool
type MapUserTypeRoutes map[string][]string
type MapPermissionRoutes map[string][]string
type MapScopeRoutes map[string][]string

type Config struct {
	mapPublicRoutes     MapPublicRoutes
	mapUserTypeRoutes   MapUserTypeRoutes
	mapPermissionRoutes MapPermissionRoutes
	mapScopeRoutes      MapScopeRoutes
}

type ConfigFunc func(c *Config)

func PublicRoutes(r MapPublicRoutes) ConfigFunc {
	return func(c *Config) {
		c.mapPublicRoutes = r
	}
}

func UserTypeRoutes(r MapUserTypeRoutes) ConfigFunc {
	return func(c *Config) {
		c.mapUserTypeRoutes = r
	}
}

func PermissionRoutes(r MapPermissionRoutes) ConfigFunc {
	return func(c *Config) {
		c.mapPermissionRoutes = r
	}
}

func ScopeRoutes(r MapScopeRoutes) ConfigFunc {
	return func(c *Config) {
		c.mapScopeRoutes = r
	}
}

func generateConfig(args ...ConfigFunc) *Config {
	c := &Config{}
	for i := range args {
		args[i](c)
	}
	return c
}
