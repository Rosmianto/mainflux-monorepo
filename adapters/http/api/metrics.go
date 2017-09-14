package api

import (
	"time"

	"github.com/go-kit/kit/metrics"
	adapter "github.com/mainflux/http-adapter"
	writer "github.com/mainflux/message-writer"
)

var _ adapter.Service = (*metricService)(nil)

type metricService struct {
	counter metrics.Counter
	latency metrics.Histogram
	adapter.Service
}

// NewMetricService instruments adapter by tracking request count and latency.
func NewMetricService(counter metrics.Counter, latency metrics.Histogram, s adapter.Service) adapter.Service {
	return &metricService{
		counter: counter,
		latency: latency,
		Service: s,
	}
}

func (ms *metricService) Send(msgs []writer.Message) {
	defer func(begin time.Time) {
		ms.counter.With("method", "send").Add(1)
		ms.latency.With("method", "send").Observe(time.Since(begin).Seconds())
	}(time.Now())

	ms.Service.Send(msgs)
}
