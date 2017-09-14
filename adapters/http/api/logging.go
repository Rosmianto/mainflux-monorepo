package api

import (
	"time"

	"github.com/go-kit/kit/log"
	adapter "github.com/mainflux/http-adapter"
	writer "github.com/mainflux/message-writer"
)

var _ adapter.Service = (*loggingService)(nil)

type loggingService struct {
	logger log.Logger
	adapter.Service
}

// NewLoggingService adds logging facilities to the adapter.
func NewLoggingService(logger log.Logger, s adapter.Service) adapter.Service {
	return &loggingService{logger, s}
}

func (ls *loggingService) Send(msgs []writer.Message) {
	defer func(begin time.Time) {
		ls.logger.Log(
			"method", "send",
			"took", time.Since(begin),
		)
	}(time.Now())

	ls.Service.Send(msgs)
}
