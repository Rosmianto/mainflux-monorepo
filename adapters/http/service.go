package adapter

import writer "github.com/mainflux/message-writer"

// Service specifies an API that must be fullfiled by the domain service
// implementation, and all of its decorators (e.g. logging & metrics).
type Service interface {
	Send([]writer.Message)
}
