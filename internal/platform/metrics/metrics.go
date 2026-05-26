package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

const (
	ResultSuccess         = "success"
	ResultValidationError = "validation_error"
	ResultUsecaseError    = "usecase_error"
)

type Metrics struct {
	Alert AlertMetrics
}

type AlertMetrics struct {
	IngestTotal *prometheus.CounterVec
}

func New() *Metrics {
	return &Metrics{
		Alert: AlertMetrics{
			IngestTotal: prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Namespace: "pulseops",
					Subsystem: "alert",
					Name:      "ingest_total",
					Help:      "Total number of alert ingest attempts.",
				},
				[]string{"result"},
			),
		},
	}
}

func (m *Metrics) Register(reg prometheus.Registerer) {
	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		m.Alert.IngestTotal,
	)
}
