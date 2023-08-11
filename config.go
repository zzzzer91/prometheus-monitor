package monitor

type ConfigOption func(c *config)

type config struct {
	path string
}

func newHttpConfig() *config {
	return &config{
		path: "/metrics",
	}
}

func (c *config) apply(opts ...ConfigOption) {
	for _, o := range opts {
		o(c)
	}
}

func WithPath(path string) ConfigOption {
	return func(c *config) {
		c.path = path
	}
}
