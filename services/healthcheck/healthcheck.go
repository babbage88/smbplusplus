// Package healthcheck provides database health check services,
// including read and insert health check functionalities.
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

// DbHeathCheckResponse represents the response structure for a database health check.
//
// swagger:model DbHeathCheckResponse
type DbHeathCheckResponse struct {
	Error     error     `json:"error"`     // Error encountered during the health check, if any.
	Status    string    `json:"status"`    // Status of the health check (e.g., "healthy", "unhealthy").
	CheckType string    `json:"checkType"` // Type of health check performed (e.g., "Read", "Create").
	Id        uuid.UUID `json:"id"`        // Unique identifier for the health check.
}

// swagger:parameters idOfdbHealthCheck
type DbHealthCheckRequest struct {
	//Type of DB HealthCheck
	// In: path
	TYPE string
}

// HealthCheckServicePgxImpl implements database health checks using pgx.
//
// This struct maintains a connection pool to execute health check queries.
type HealthCheckServicePgxImpl struct {
	DbConn *pgxpool.Pool // Database connection pool.
}

// HcDbParser defines methods for parsing database health check results.
type HcDbParser interface {
	ParseDbReadHealthCheck(db smbplusplus_db.DbHealthCheckReadRow)
	ParseDbHealthCheck(db smbplusplus_db.HealthCheck)
}

// HealthCheckService defines the interface for health check operations.
type HealthCheckService interface {
	DbReadHealthCheck() DbHeathCheckResponse
	DbInsertHealthCheck() DbHeathCheckResponse
	DbDeleteHealthCheck(id uuid.UUID) DbHeathCheckResponse
	InsertAndDeleteHealthCheck() DbHeathCheckResponse
	DbHealthCheckHandler() http.Handler
}

// ParseDbReadHealthCheck parses a database read health check result.
//
// It updates the DbHeathCheckResponse struct based on the query result.
func (dbhc *DbHeathCheckResponse) ParseDbReadHealthCheck(db smbplusplus_db.DbHealthCheckReadRow) {
	dbhc.CheckType = db.CheckType.String
	dbhc.Id = db.ID
	dbhc.Status = db.Status.String
	if db.Status.String != "healthy" {
		dbhc.Error = fmt.Errorf("database did not respond with healthy status")
	} else {
		dbhc.Error = nil
	}
}

// ParseDbHealthCheck parses a general database health check result.
//
// It updates the DbHeathCheckResponse struct based on the query result.
func (dbhc *DbHeathCheckResponse) ParseDbHealthCheck(db smbplusplus_db.HealthCheck) {
	dbhc.CheckType = db.CheckType.String
	dbhc.Id = db.ID
	dbhc.Status = db.Status.String
	if db.Status.String != "healthy" {
		dbhc.Error = fmt.Errorf("database did not respond with healthy status")
	} else {
		dbhc.Error = nil
	}
}

// DbReadHealthCheck performs a database read health check.
//
// It queries the database and returns a response indicating the health status.
func (h *HealthCheckServicePgxImpl) DbReadHealthCheck() DbHeathCheckResponse {
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

// DbInsertHealthCheck performs a database insert health check.
//
// It executes an insert operation and returns a response indicating the health status.
func (h *HealthCheckServicePgxImpl) DbInsertHealthCheck() DbHeathCheckResponse {
	dbHealth := &DbHeathCheckResponse{CheckType: "Create"}
	queries := smbplusplus_db.New(h.DbConn)
	qry, err := queries.DbHealthCheckInsert(context.Background())
	if err != nil {
		slog.Error("Error executing DbInsertHealthCheck query", slog.String("error", err.Error()))
		dbHealth.Error = err
		return *dbHealth
	}
	dbHealth.ParseDbHealthCheck(qry)

	return *dbHealth
}

func (h *HealthCheckServicePgxImpl) DbDeleteHealthCheck(id uuid.UUID) DbHeathCheckResponse {
	dbHealth := &DbHeathCheckResponse{CheckType: "Delete"}
	queries := smbplusplus_db.New(h.DbConn)
	qry, err := queries.DbHealthCheckDelete(context.Background(), id)
	if err != nil {
		slog.Error("Error executing DbInsertHealthCheck query", slog.String("error", err.Error()))
		dbHealth.Error = err
		return *dbHealth
	}
	dbHealth.ParseDbHealthCheck(qry)

	return *dbHealth
}

func (h *HealthCheckServicePgxImpl) InsertAndDeleteHealthCheck() DbHeathCheckResponse {
	dbHealth := h.DbInsertHealthCheck()
	delRecord := h.DbDeleteHealthCheck(dbHealth.Id)
	return delRecord
}
