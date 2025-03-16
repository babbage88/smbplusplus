// Package main smbplusplus API.
//
// Terms Of Service:
//
// there are no TOS at this moment, use at your own risk we take no responsibility
//
//		Version: v0.0.1
//		License: N/A
//		Contact: Justin Trahan<test@trahan.dev>
//
//		Consumes:
//		- application/json
//
//		Produces:
//		- application/json
//
//	    Security:
//	    - bearer:
//
//	    SecurityDefinitions:
//	      bearer:
//	         type: apiKey
//	         name: Authorization
//	         in: header
//
// swagger:meta
package main

import (
	_ "embed"
	"flag"
	"log/slog"
	"os"

	"github.com/babbage88/smbplusplus/api"
	"github.com/babbage88/smbplusplus/database/s2_pgxpool"
	"github.com/babbage88/smbplusplus/services/healthcheck"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

//go:embed swagger.yaml
var swaggerSpec []byte

func loadEnvVars(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		slog.Error("Error loading .env file", slog.String("path", path), slog.String("error", err.Error()))
	}
	return err
}

func initPgConnPool(dbUrl string) *pgxpool.Pool {
	connPool := s2_pgxpool.PgPoolInit(dbUrl)
	return connPool
}

func main() {
	var dbUrl string
	var envFile string = ".env"

	srvport := flag.String("srvadr", ":8559", "Address and port that http server will listed on. :8559 is default")
	flag.StringVar(&dbUrl, "db", "", "Overide for the database connection url, otherwise DATABASE_URL env var will be used.")
	flag.StringVar(&envFile, "env-file", ".env", "Env file for loading environment variables.")
	flag.Parse()
	loadEnvVars(envFile)
	if dbUrl == "" {
		dbUrl = os.Getenv("DATABASE_URL")
	}

	dbConn := initPgConnPool(dbUrl)
	api.SetSwaggerSpec(swaggerSpec)
	healthCheckService := healthcheck.HealthCheckService{DbConn: dbConn}
	err := api.StartApiServer(srvport, &healthCheckService)
	if err != nil {
		slog.Error("error creating new server instance", slog.String("error", err.Error()))
	}
}
