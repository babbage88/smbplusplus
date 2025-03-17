package healthcheck

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type DbHealthCheckHandler struct {
	service *HealthCheckService
}

// swagger:route GET /health/db/{TYPE} dbHealthCheck idOfdbHealthCheck
// Performs database health check and returns a respoonse. Currently defaults to Read, but takes the type (eg: read, write, update, insert)
// as a url path parameter
//
// security:
// - bearer:
// responses:
//   200: DbHeathCheckResponse

func (h *DbHealthCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	checkType := r.PathValue("TYPE")
	if checkType == "insert" || checkType == "write" {
		slog.Error("insert check not yet implemented")
		http.Error(w, "write check not implemented", http.StatusNotFound)
		return
	}

	readCheck := h.service.GetDbReadHealthCheck()
	if readCheck.Error != nil {
		slog.Error("Error running db read healthcheck", slog.String("error", readCheck.Error.Error()))
		http.Error(w, "Failed to run database healthcheck query: "+readCheck.Error.Error(), http.StatusInternalServerError)
		return
	}

	hcResponse, err := json.Marshal(readCheck)
	if err != nil {
		slog.Error("Error marshaling response", slog.String("error", err.Error()))
		http.Error(w, "error marshaling response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(hcResponse)
	slog.Info("Response sent successfully")
}

func (h *HealthCheckService) DbHealthCheckHandler() http.Handler {
	return &DbHealthCheckHandler{service: h}
}
