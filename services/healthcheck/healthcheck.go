package healthcheck

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	smbplusplus_db "github.com/babbage88/smbplusplus/database/smbplusplus_pg"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DbHeathCheckResponse struct {
	Error     error     `json:"error"`
	Status    string    `json:"status"`
	CheckType string    `json:"checkType"`
	Id        uuid.UUID `json:"id"`
}

type HealthCheckService struct {
	DbConn *pgxpool.Pool
}

type IHealthCheckService interface {
	GetDbReadHealthCheck() DbHeathCheckResponse
	DbReadHealthCheckHandler() func(http.ResponseWriter, *http.Request)
	ParseDbReadHealthCheck(db smbplusplus_db.DbHealthCheckReadRow)
}

func (dbhc *DbHeathCheckResponse) ParseDbReadHealthCheck(db smbplusplus_db.DbHealthCheckReadRow) {
	dbhc.CheckType = db.CheckType.String
	dbhc.Id = db.ID
	dbhc.Status = db.Status.String
	if db.Status.String != "healthy" {
		dbhc.Error = fmt.Errorf("database did not responde with healthy status.")
	} else {
		dbhc.Error = nil
	}
}

func (h *HealthCheckService) GetDbReadHealthCheck() DbHeathCheckResponse {
	dbHealth := &DbHeathCheckResponse{CheckType: "Read"}
	queries := smbplusplus_db.New(h.DbConn)
	qry, err := queries.DbHealthCheckRead(context.Background())
	if err != nil {
		slog.Error("Error executing DbReadHealthCheck query", slog.String("error", err.Error()))
		dbHealth.Error = err
		return *dbHealth
	}
	dbHealth.ParseDbReadHealthCheck(qry)

	return *dbHealth
}

// swagger:route GET /health/db/{TYPE} dbHealthCheck idOfdbHealthCheck
// Performs database health check and returns a respoonse. Currently defaults to Read, but takes the type (eg: read, write, update, insert)
// as a url path parameter
//
// security:
// - bearer:
// responses:
//   200: GetUserByIdResponse

func (h *HealthCheckService) DbHealthCheckHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		checkType := r.PathValue("TYPE")
		if checkType == "insert" || checkType == "write" {
			slog.Error("insert check not yet implemented")
			http.Error(w, "write check not implemented", http.StatusNotFound)
			return
		} else {
			readCheck := h.GetDbReadHealthCheck()

			if readCheck.Error != nil {
				slog.Error("Error running db read healthcheck", slog.String("error", readCheck.Error.Error()))
				http.Error(w, "Failed to run database healthcheck query: "+readCheck.Error.Error(), http.StatusInternalServerError)
				return
			}

			hcResponse, err := json.Marshal(readCheck)
			if err != nil {
				slog.Error("Error marshing response", slog.String("error", err.Error()))
				http.Error(w, "error marshaling response"+err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(hcResponse)
			slog.Info("Response sent successfully")
		}
	}
}
