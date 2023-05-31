package monitor

type httpConfig struct {
	path string
}

func newHttpConfig() *httpConfig {
	return &httpConfig{
		path: "/metrics",
	}
}

func (c *httpConfig) apply(opts ...httpConfigOption) {
	for _, o := range opts {
		o(c)
	}
}

type httpConfigOption func(c *httpConfig)

func WithPath(path string) httpConfigOption {
	return func(c *httpConfig) {
		c.path = path
	}
}
