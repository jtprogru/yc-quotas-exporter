package yandexcloud

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/jtprogru/yc-quotas-exporter/internal/config"
)

type Exporter struct {
	Config *config.Config
	Client *Client
}

func NewExporter(cfg *config.Config) (*Exporter, error) {
	client, err := NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &Exporter{
		Config: cfg,
		Client: client,
	}, nil
}

func (e *Exporter) Run() error {
	// Create a new QuotaExporter.
	exporter, err := NewQuotaExporter(e.Config)
	if err != nil {
		return err
	}

	// Export metrics.
	go func() {
		for {
			err := exporter.ExportMetrics()
			if err != nil {
				log.Printf("Error ExportMetrics: %v", err)
			}
			time.Sleep(time.Duration(e.Config.Timeout) * time.Second)
		}
	}()

	// Expose the registered metrics via HTTP with logging middleware.
	http.Handle("/metrics", loggingMiddleware(promhttp.Handler()))

	// Start the HTTP server with timeout.
	srv := &http.Server{
		Addr:         net.JoinHostPort(e.Config.Host, e.Config.Port),
		ReadTimeout:  time.Duration(e.Config.Timeout) * time.Second,
		WriteTimeout: time.Duration(e.Config.Timeout) * time.Second,
	}
	log.Printf("Starting server on %s", srv.Addr)
	return srv.ListenAndServe()
}

// loggingMiddleware is a middleware handler that logs the incoming HTTP requests.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		log.Printf("Completed %s in %v", r.URL.Path, time.Since(start))
	})
}
