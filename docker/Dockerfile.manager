FROM golang:1.8-alpine AS builder
WORKDIR /go/src/github.com/mainflux/manager
COPY . .
RUN cd cmd && CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo -o app

FROM scratch
COPY --from=builder /go/src/github.com/mainflux/manager/cmd/app /
EXPOSE 8180
ENTRYPOINT ["/app"]
