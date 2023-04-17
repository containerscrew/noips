package infra

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type MetricsData struct {
	subnet_available_ipv4 *prometheus.GaugeVec
}

type PrometheusInputData struct {
	NewRegistry prometheus.Registry
	Handler     http.Handler
}

func NewMetrics() *PrometheusInputData {
	return &PrometheusInputData{
		NewRegistry: *prometheus.NewRegistry(),
	}
}

func (p *PrometheusInputData) RegisterMetrics() {
	m := &MetricsData{
		subnet_available_ipv4: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "aws_subnets_free_ipv4",
			Name:      "info",
			Help:      "Information about the my aws subnets.",
		},
			[]string{"subnet_id", "subnet_cidr"}),
	}

	p.NewRegistry.MustRegister(m.subnet_available_ipv4)
}

func (p *PrometheusInputData) RegisterHandler() http.Handler {
	return promhttp.HandlerFor(&p.NewRegistry, promhttp.HandlerOpts{})
}

func (p *PrometheusInputData) UpdateMetrics() {
	panic(fmt.Errorf("implementing function"))
}
