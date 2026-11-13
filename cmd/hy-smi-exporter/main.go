// Command hy-smi-exporter is a Prometheus exporter for Hygon DCU
// (Data Center Unit) / HCU accelerators. It shells out to the
// read-only `hy-smi` command and exposes the parsed
// per-device metrics under the `hy_smi_` prefix on /metrics.
//
// Only read-only hy-smi options are ever invoked. No --set*/,
// --rasinject, or --load*driver write operations are issued.
//
// The exporter is a single static Go binary with no cgo and no
// runtime dependencies beyond the hy-smi binary itself, which must
// be resolvable via $PATH or pointed at by --hy-smi.path.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/di-matrix-ai/hy-smi-exporter/internal/config"
	"github.com/di-matrix-ai/hy-smi-exporter/internal/exporter"
	"github.com/di-matrix-ai/hy-smi-exporter/internal/hysmi"
)

// Version metadata is injected at build time via -ldflags.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cfg := config.FromEnv()

	var (
		listenAddr            = flag.String("web.listen-address", cfg.ListenAddr, "Address to listen on for /metrics.")
		metricsPath           = flag.String("web.telemetry-path", cfg.MetricsPath, "Path under which metrics are exposed.")
		logLevel              = flag.String("log.level", cfg.LogLevel, "Log level: debug, info, warn, error.")
		binaryPath            = flag.String("hy-smi.path", cfg.BinaryPath, "Path to the hy-smi binary. Resolved through $PATH if unqualified.")
		timeout               = flag.Duration("hy-smi.timeout", cfg.Timeout, "Timeout for a single hy-smi invocation.")
		collectOneSecondStats = flag.Bool("hy-smi.collect-1s-stats", cfg.CollectOneSecondStats, "Collect detailed one-second utilization and bandwidth statistics.")
		showVersion           = flag.Bool("version", false, "Print version information and exit.")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("hy-smi-exporter version=%s commit=%s built=%s\n", version, commit, date)
		return
	}

	logger := newLogger(*logLevel)
	slog.SetDefault(logger)

	// Re-resolve config from flags so the collector sees the final
	// values.
	finalCfg := config.Config{
		ListenAddr:            *listenAddr,
		MetricsPath:           *metricsPath,
		LogLevel:              *logLevel,
		BinaryPath:            *binaryPath,
		Timeout:               *timeout,
		CollectOneSecondStats: *collectOneSecondStats,
	}

	registry := prometheus.NewRegistry()
	registry.MustRegister(
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewGoCollector(),
		exporter.NewCollector(hysmi.New(finalCfg.BinaryPath, finalCfg.Timeout, finalCfg.CollectOneSecondStats), logger),
	)

	mux := http.NewServeMux()
	mux.Handle(*metricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprintf(w, `<!doctype html>
<html><head><title>hy-smi-exporter</title></head><body>
<h1>hy-smi-exporter</h1>
<p>version: %s</p>
<ul>
  <li><a href="%s">%s</a></li>
  <li><a href="/healthz">/healthz</a></li>
</ul>
</body></html>
`, version, *metricsPath, *metricsPath)
	})

	srv := &http.Server{
		Addr:              *listenAddr,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Info("hy-smi-exporter listening",
			"addr", *listenAddr,
			"binary", finalCfg.BinaryPath,
			"timeout", finalCfg.Timeout.String(),
			"collect_1s_stats", finalCfg.CollectOneSecondStats,
		)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server error", "err", err)
			stop()
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", "err", err)
		os.Exit(1)
	}
}

func newLogger(level string) *slog.Logger {
	var lvl slog.Level
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		lvl = slog.LevelDebug
	case "warn", "warning":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}
	return slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: lvl}))
}
