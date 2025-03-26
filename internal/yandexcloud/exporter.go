package yandexcloud

import (
	"log"
	"net/http"
	"time"

	"github.com/jtprogru/yc-quotas-exporter/internal/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

func (e *Exporter) Run() {

	// Create a new QuotaExporter
	exporter, err := NewQuotaExporter(e.Config)
	if err != nil {
		log.Fatalf("Error creating QuotaExporter: %v", err)
	}

	// Export metrics
	go exporter.ExportMetrics()

	// Expose the registered metrics via HTTP with logging middleware
	http.Handle("/metrics", loggingMiddleware(promhttp.Handler()))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// loggingMiddleware is a middleware handler that logs the incoming HTTP requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		log.Printf("Completed %s in %v", r.URL.Path, time.Since(start))
	})
}
