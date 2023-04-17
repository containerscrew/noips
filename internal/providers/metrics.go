package providers

import "net/http"

type Metrics interface {
	RegisterMetrics()
	RegisterHandler() http.Handler
	UpdateMetrics()
}

/*
		var (
		subnetId   = "subnet_id"
		subnetCIDR = "subnet_cidr"
	)

	type metrics struct {
		info *prometheus.GaugeVec
	}

	func NewMetrics(reg prometheus.Registerer) *metrics {
		m := &metrics{
			info: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: "aws_subnets_free_ipv4",
				Name:      "info",
				Help:      "Information about the my aws subnets.",
			},
				[]string{subnetId, subnetCIDR}),
		}
		reg.MustRegister(m.info)
		return m
	}

	//Prometheus monitoring
	reg := prometheus.NewRegistry()
	m := NewMetrics(reg)

	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
*/
