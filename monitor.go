package monitor

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Monitor is an object that uses to set gin server monitor.
type Monitor struct {
	metrics  map[string]*Metric
	registry *prometheus.Registry

	httpHandler     http.Handler
	httpHandlerOnce sync.Once
}

func NewMonitor(registry *prometheus.Registry) *Monitor {
	if registry == nil {
		registry = prometheus.DefaultRegisterer.(*prometheus.Registry)
	}
	return &Monitor{
		metrics:  make(map[string]*Metric),
		registry: registry,
	}
}

// GetRegistry used to get prometheus registry.
func (m *Monitor) GetRegistry() *prometheus.Registry {
	return m.registry
}

// GetMetric used to get metric object by metric_name.
func (m *Monitor) GetMetric(name string) *Metric {
	return m.metrics[name]
}

// AddMetric add custom monitor metric.
func (m *Monitor) AddMetric(metric *Metric) {
	if _, ok := m.metrics[metric.Name]; ok {
		panic(fmt.Sprintf("metric '%s' is existed", metric.Name))
	}
	if metric.Name == "" {
		panic("metric name cannot be empty.")
	}
	if f, ok := promTypeHandler[metric.Type]; ok {
		f(metric)
		m.registry.MustRegister(metric.vec)
		m.metrics[metric.Name] = metric
	} else {
		panic(fmt.Sprintf("metric type '%d' not existed.", metric.Type))
	}
}

// HttpHandler used to generate http handler.
// Usage: `http.Handle(path, m.HttpHandler())`
func (m *Monitor) HttpHandler() http.Handler {
	m.httpHandlerOnce.Do(func() {
		m.httpHandler = promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{ErrorHandling: promhttp.ContinueOnError})
	})
	return m.httpHandler
}
