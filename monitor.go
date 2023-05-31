package monitor

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

// Monitor is an object that uses to set gin server monitor.
type Monitor struct {
	*prometheus.Registry

	metrics map[string]*Metric
}

func NewMonitor(registry *prometheus.Registry) *Monitor {
	if registry == nil {
		panic("registry can not be nil")
	}
	return &Monitor{
		metrics:  make(map[string]*Metric),
		Registry: registry,
	}
}

var (
	m = NewMonitor(prometheus.DefaultRegisterer.(*prometheus.Registry))
)

func DefaultMonitor() *Monitor {
	return m
}

// GetRegistry used to get prometheus registry.
func (m *Monitor) GetRegistry() *prometheus.Registry {
	return m.Registry
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
		m.MustRegister(metric.vec)
		m.metrics[metric.Name] = metric
	} else {
		panic(fmt.Sprintf("metric type '%d' not existed.", metric.Type))
	}
}
