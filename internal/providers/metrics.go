package providers

type Metrics interface {
	RegisterMetrics()
	UpdateMetrics()
}
