package app

import (
	"net/http"
)

type Metrics interface {
	RegisterMetrics()
	RegisterHandler() http.Handler
	UpdateMetrics(subnet SubnetsData)
}
