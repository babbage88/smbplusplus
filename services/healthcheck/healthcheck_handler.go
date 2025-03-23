package healthcheck

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
)

// DbHealthCheckHandler handles HTTP requests for database health checks.
type DbHealthCheckHandler struct {
	service *HealthCheckServicePgxImpl
}

// swagger:route GET /health/db/{TYPE} dbHealthCheck idOfdbHealthCheck
//
// Performs database health check and returns a respoonse.
// Determines the type of health check to perform based on the URL path parameter "type".
// If the check type corresponds to an insert/write operation, it executes an insert health check.
// Otherwise, it defaults to a read health check.
//
// Supported check types for insert operations: "insert", "delete", "read".
// All other types default to a read health check.
//
// security:
// - bearer:
// responses:
//   200: DbHeathCheckResponse

func (h *DbHealthCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	pathParamVal := r.PathValue("type")
	checkType := strings.ToLower(pathParamVal)
	if checkType == "insert" || checkType == "write" || checkType == "create" || checkType == "new" {
		insertCheck := h.service.DbInsertHealthCheck()
		defer h.service.DbDeleteHealthCheck(insertCheck.Id)
		insertResponse, err := json.Marshal(insertCheck)
		if err != nil {
			slog.Error("Error marshaling response", slog.String("error", err.Error()))
			http.Error(w, "error marshaling response: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(insertResponse)
		slog.Info("Response sent successfully")
		return

	}

	if checkType == "delete" || checkType == "drop" || checkType == "rm" || checkType == "remove" {
		deleteDbHc := h.service.InsertAndDeleteHealthCheck()
		delResponse, err := json.Marshal(deleteDbHc)
		if err != nil {
			slog.Error("Error marshaling response", slog.String("error", err.Error()))
			http.Error(w, "error marshaling response: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(delResponse)
		slog.Info("Response sent successfully")
		return

	}

	readCheck := h.service.DbReadHealthCheck()
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

// DbHealthCheckHandler wraps a HealthCheckServicePgxImpl to implement http.Handler interface
//
// This handler listens for Database Health checks and performs the corresponding test.
func (h *HealthCheckServicePgxImpl) DbHealthCheckHandler() http.Handler {
	return &DbHealthCheckHandler{service: h}
}
