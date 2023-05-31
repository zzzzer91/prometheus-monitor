package monitor

type httpConfig struct {
	port uint16
	path string
}

func newHttpConfig() *httpConfig {
	return &httpConfig{
		port: 9100,
		path: "/metrics",
	}
}

func (c *httpConfig) apply(opts ...httpConfigOption) {
	for _, o := range opts {
		o(c)
	}
}

type httpConfigOption func(c *httpConfig)

func WithPort(port uint16) httpConfigOption {
	return func(c *httpConfig) {
		c.port = port
	}
}

func WithPath(path string) httpConfigOption {
	return func(c *httpConfig) {
		c.path = path
	}
}
