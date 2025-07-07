FROM golang:alpine AS builder
WORKDIR /build
ADD go.mod .
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o sniffer main.go
FROM alpine
WORKDIR /build
COPY --from=builder /build/sniffer /build/sniffer
RUN chmod +x sniffer
CMD ["./sniffer"]