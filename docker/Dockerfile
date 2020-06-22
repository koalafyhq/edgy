FROM golang:alpine as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build

# TODO(@faultable): fix this
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o edgy ./cmd/main.go

FROM alpine:latest
COPY --from=builder /build/edgy /app/
WORKDIR /app

RUN apk --no-cache add ca-certificates curl

EXPOSE 3000

CMD ["./edgy"]
