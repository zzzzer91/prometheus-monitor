package monitor

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

type MetricType int

const (
	None MetricType = iota
	Counter
	Gauge
	Histogram
	Summary
)

// Metric defines a metric object. Users can use it to save
// metric data. Every metric should be globally unique by name.
type Metric struct {
	Type        MetricType
	NameSpace   string
	Name        string
	Description string
	Labels      []string
	Buckets     []float64
	Objectives  map[float64]float64

	vec prometheus.Collector
}

// SetGaugeValue set data for Gauge type Metric.
func (m *Metric) SetGaugeValue(labelValues []string, value float64) error {
	if m.Type == None {
		return errors.Errorf("metric '%s' not existed.", m.Name)
	}

	if m.Type != Gauge {
		return errors.Errorf("metric '%s' not Gauge type", m.Name)
	}
	m.vec.(*prometheus.GaugeVec).WithLabelValues(labelValues...).Set(value)
	return nil
}

// Inc increases value for Counter/Gauge type metric, increments
// the counter by 1
func (m *Metric) Inc(labelValues []string) error {
	if m.Type == None {
		return errors.Errorf("metric '%s' not existed.", m.Name)
	}

	if m.Type != Gauge && m.Type != Counter {
		return errors.Errorf("metric '%s' not Gauge or Counter type", m.Name)
	}
	switch m.Type {
	case Counter:
		m.vec.(*prometheus.CounterVec).WithLabelValues(labelValues...).Inc()
	case Gauge:
		m.vec.(*prometheus.GaugeVec).WithLabelValues(labelValues...).Inc()
	}
	return nil
}

// Add adds the given value to the Metric object. Only
// for Counter/Gauge type metric.
func (m *Metric) Add(labelValues []string, value float64) error {
	if m.Type == None {
		return errors.Errorf("metric '%s' not existed.", m.Name)
	}

	if m.Type != Gauge && m.Type != Counter {
		return errors.Errorf("metric '%s' not Gauge or Counter type", m.Name)
	}
	switch m.Type {
	case Counter:
		m.vec.(*prometheus.CounterVec).WithLabelValues(labelValues...).Add(value)
	case Gauge:
		m.vec.(*prometheus.GaugeVec).WithLabelValues(labelValues...).Add(value)
	}
	return nil
}

// Observe is used by Histogram and Summary type metric to
// add observations.
func (m *Metric) Observe(labelValues []string, value float64) error {
	if m.Type == 0 {
		return errors.Errorf("metric '%s' not existed.", m.Name)
	}
	if m.Type != Histogram && m.Type != Summary {
		return errors.Errorf("metric '%s' not Histogram or Summary type", m.Name)
	}
	switch m.Type {
	case Histogram:
		m.vec.(*prometheus.HistogramVec).WithLabelValues(labelValues...).Observe(value)
	case Summary:
		m.vec.(*prometheus.SummaryVec).WithLabelValues(labelValues...).Observe(value)
	}
	return nil
}

var (
	promTypeHandler = map[MetricType]func(metric *Metric){
		Counter:   counterHandler,
		Gauge:     gaugeHandler,
		Histogram: histogramHandler,
		Summary:   summaryHandler,
	}
)

func counterHandler(metric *Metric) {
	metric.vec = prometheus.NewCounterVec(
		prometheus.CounterOpts{Namespace: metric.NameSpace, Name: metric.Name, Help: metric.Description},
		metric.Labels,
	)
}

func gaugeHandler(metric *Metric) {
	metric.vec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{Namespace: metric.NameSpace, Name: metric.Name, Help: metric.Description},
		metric.Labels,
	)
}

func histogramHandler(metric *Metric) {
	if len(metric.Buckets) == 0 {
		panic(fmt.Sprintf("metric '%s' is histogram type, cannot lose bucket param.", metric.Name))
	}
	metric.vec = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{Namespace: metric.NameSpace, Name: metric.Name, Help: metric.Description, Buckets: metric.Buckets},
		metric.Labels,
	)
}

func summaryHandler(metric *Metric) {
	if len(metric.Objectives) == 0 {
		panic(fmt.Sprintf("metric '%s' is summary type, cannot lose objectives param.", metric.Name))
	}
	prometheus.NewSummaryVec(
		prometheus.SummaryOpts{Namespace: metric.NameSpace, Name: metric.Name, Help: metric.Description, Objectives: metric.Objectives},
		metric.Labels,
	)
}
