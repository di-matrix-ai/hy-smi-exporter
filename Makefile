namespace      := hy_smi
binary         := hy-smi-exporter
ldflags        := -X main.version=$(shell git describe --tags --always --dirty 2>/dev/null || echo dev) \
                  -X main.commit=$(shell git rev-parse --short HEAD 2>/dev/null || echo none) \
                  -X main.date=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)

GOFLAGS        := -trimpath -ldflags "$(ldflags)"

.PHONY: build test lint fmt vet release-check clean

build:
	go build $(GOFLAGS) -o bin/$(binary) ./cmd/$(binary)

test:
	go test -v ./...

fmt:
	gofmt -s -w .

vet:
	go vet ./...

lint: vet

release-check:
	EXPECTED_MODULE=github.com/di-matrix-ai/hy-smi-exporter ./scripts/check-release-materials.sh

clean:
	rm -rf bin/
