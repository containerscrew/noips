package infra

import (
	"fmt"
	"github.com/nanih98/noips/pkg/app"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

type MetricsData struct {
	subnet_available_ipv4 *prometheus.GaugeVec
}

type PrometheusInputData struct {
	NewRegistry prometheus.Registry
	Handler     http.Handler
	Metrics     *MetricsData
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
	p.Metrics = m
}

func (p *PrometheusInputData) RegisterHandler() http.Handler {
	return promhttp.HandlerFor(&p.NewRegistry, promhttp.HandlerOpts{})
}

func (p *PrometheusInputData) UpdateMetrics(subnet app.SubnetsData) {
	log.Println("Refreshing data for", subnet.ID, subnet.CIDR, subnet.AvailableIPV4)
	fmt.Println(subnet.ID, subnet.CIDR, subnet.AvailableIPV4)
	p.Metrics.subnet_available_ipv4.With(prometheus.Labels{"subnet_id": subnet.ID, "subnet_cidr": subnet.CIDR}).Set(float64(subnet.AvailableIPV4))
}
