package api

import (
	_ "embed"
	"log/slog"
	"net/http"

	"github.com/babbage88/smbplusplus/internal/cors"
	"github.com/babbage88/smbplusplus/internal/swaggerui"
	"github.com/babbage88/smbplusplus/services/healthcheck"
	"github.com/babbage88/smbplusplus/services/s2auth"
	"github.com/babbage88/smbplusplus/services/s2usercrud"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var SwaggerSpec []byte

func SetSwaggerSpec(swaggerSpec []byte) {
	SwaggerSpec = swaggerSpec
}

func StartApiServer(srvadr *string, hc healthcheck.HealthCheckService, auth_svc s2auth.AuthService, userCrud_svc s2usercrud.UserCRUD) error {
	mux := http.NewServeMux()
	mux.Handle("GET /health/db/{type}", cors.CORSWithGET(hc.DbHealthCheckHandler()))
	mux.Handle("POST /login", cors.CORSWithPOST(http.HandlerFunc(s2auth.LoginHandleFunc(auth_svc))))
	mux.Handle("POST /create/user", cors.CORSWithPOST(http.HandlerFunc(s2usercrud.CreateUserHandler(userCrud_svc))))

	mux.Handle("/metrics", promhttp.Handler())
	// Add Swagger UI handler
	mux.Handle("/swaggerui/", http.StripPrefix("/swaggerui", swaggerui.ServeSwaggerUI(SwaggerSpec)))
	err := http.ListenAndServe(*srvadr, mux)
	if err != nil {
		slog.Error("Failed to start server", slog.String("Error", err.Error()))
	}
	return err
}
