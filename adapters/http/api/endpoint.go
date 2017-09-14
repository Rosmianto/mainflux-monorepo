package api

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	adapter "github.com/mainflux/http-adapter"
	writer "github.com/mainflux/message-writer"
)

func sendMessageEndpoint(svc adapter.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		messages := request.([]writer.Message)
		svc.Send(messages)
		return nil, nil
	}
}
