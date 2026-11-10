FROM golang:1.23-alpine AS builder
WORKDIR /src
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags "-s -w \
    -X main.version=$(git describe --tags --always 2>/dev/null || echo dev) \
    -X main.commit=$(git rev-parse --short HEAD 2>/dev/null || echo none)" \
    -o /out/hy-smi-exporter ./cmd/hy-smi-exporter

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=builder /out/hy-smi-exporter /usr/local/bin/hy-smi-exporter
EXPOSE 9393
USER 65534:65534
ENTRYPOINT ["/usr/local/bin/hy-smi-exporter"]
