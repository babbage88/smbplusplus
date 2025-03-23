package healthcheck

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	smbplusplus_db "github.com/babbage88/smbplusplus/database/smbplusplus_pg"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// swagger:model DbHeathCheckResponse
type DbHeathCheckResponse struct {
	Error     error     `json:"error"`
	Status    string    `json:"status"`
	CheckType string    `json:"checkType"`
	Id        uuid.UUID `json:"id"`
}

// swagger:parameters idOfdbHealthCheck
type DbHealthCheckRequest struct {
	//Type of DB HealthCheck
	// In: path
	TYPE string
}

type HealthCheckService struct {
	DbConn *pgxpool.Pool
}

type IHealthCheckService interface {
	GetDbReadHealthCheck() DbHeathCheckResponse
	DbReadHealthCheckHandler() func(http.ResponseWriter, *http.Request)
	ParseDbReadHealthCheck(db smbplusplus_db.DbHealthCheckReadRow)
	ParseDbHealthCheck(db smbplusplus_db.HealthCheck)
}

func (dbhc *DbHeathCheckResponse) ParseDbReadHealthCheck(db smbplusplus_db.DbHealthCheckReadRow) {
	dbhc.CheckType = db.CheckType.String
	dbhc.Id = db.ID
	dbhc.Status = db.Status.String
	if db.Status.String != "healthy" {
		dbhc.Error = fmt.Errorf("database did not responde with healthy status")
	} else {
		dbhc.Error = nil
	}
}

func (dbhc *DbHeathCheckResponse) ParseDbHealthCheck(db smbplusplus_db.HealthCheck) {
	dbhc.CheckType = db.CheckType.String
	dbhc.Id = db.ID
	dbhc.Status = db.Status.String
	if db.Status.String != "healthy" {
		dbhc.Error = fmt.Errorf("database did not responde with healthy status")
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

func (h *HealthCheckService) DbInsertHealthCheck() DbHeathCheckResponse {
	dbHealth := &DbHeathCheckResponse{CheckType: "Create"}
	queries := smbplusplus_db.New(h.DbConn)
	qry, err := queries.DbHealthCheckInsert(context.Background())
	if err != nil {
		slog.Error("Error executing DbReadHealthCheck query", slog.String("error", err.Error()))
		dbHealth.Error = err
		return *dbHealth
	}
	dbHealth.ParseDbHealthCheck(qry)

	return *dbHealth
}
