package api

import (
	_ "embed"
	"log/slog"
	"net/http"

	"github.com/babbage88/smbplusplus/internal/swaggerui"
	"github.com/babbage88/smbplusplus/services/healthcheck"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var SwaggerSpec []byte

func SetSwaggerSpec(swaggerSpec []byte) {
	SwaggerSpec = swaggerSpec
}

func StartApiServer(srvadr *string, hc *healthcheck.HealthCheckService) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/health/db/{TYPE}", hc.DbHealthCheckHandler())

	mux.Handle("/metrics", promhttp.Handler())
	// Add Swagger UI handler
	mux.Handle("/swaggerui/", http.StripPrefix("/swaggerui", swaggerui.ServeSwaggerUI(SwaggerSpec)))
	err := http.ListenAndServe(*srvadr, mux)
	if err != nil {
		slog.Error("Failed to start server", slog.String("Error", err.Error()))
	}
	return err
}
