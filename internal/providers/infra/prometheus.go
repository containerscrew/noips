package infra

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	subnet_available_ipv4 *prometheus.GaugeVec
}

type PrometheusInputData struct {
	NewRegistry           prometheus.Registerer
}

func NewMetrics() *PrometheusInputData {
	return &PrometheusInputData{
		NewRegistry:           prometheus.Registerer().DefaultRegisterer{},
		:
	}
}

func (p *PrometheusInputData) RegisterMetric() {
	m := &Metrics{
		subnet_available_ipv4: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "aws_subnets_free_ipv4",
			Name:      "info",
			Help:      "Information about the my aws subnets.",
		},
			[]string{"subnet_id", "subnet_cidr"}),
	}
	p.NewRegistry.MustRegister(m.subnet_available_ipv4)
}
