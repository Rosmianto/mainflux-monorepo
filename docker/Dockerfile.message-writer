FROM golang:1.8-alpine AS builder
WORKDIR /go/src/github.com/mainflux/message-writer
COPY . .
RUN cd cmd && CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo -o app

FROM scratch
COPY --from=builder /go/src/github.com/mainflux/message-writer/cmd/app /
ENTRYPOINT ["/app"]
