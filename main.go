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
	"log"
	"log/slog"
	"os"

	"github.com/babbage88/smbplusplus/api"
	"github.com/babbage88/smbplusplus/database/s2_pgxpool"
	"github.com/babbage88/smbplusplus/services/healthcheck"
	"github.com/babbage88/smbplusplus/services/s2auth"
	"github.com/babbage88/smbplusplus/services/s2usercrud"
	"github.com/google/uuid"
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
	var isDevelopment bool
	var resetDevUSer bool
	var debugProfile bool

	srvport := flag.String("srvadr", ":8995", "Address and port that http server will listed on. :8559 is default")
	flag.StringVar(&dbUrl, "db", "", "Overide for the database connection url, otherwise DATABASE_URL env var will be used.")
	flag.StringVar(&envFile, "env-file", ".env", "Env file for loading environment variables.")
	flag.BoolVar(&isDevelopment, "development", false, "Flag to start application in Development mode, env vars loaded from env-file")
	flag.BoolVar(&debugProfile, "debug", false, "Flag to enable go Debug/Profiling.")

	flag.BoolVar(&resetDevUSer, "reset-devuser", false, "Flag when set to true, the builtin Admin user (devuser) will have it's password set from the ENV variable DEV_APP_PASS")
	flag.Parse()
	if isDevelopment {
		slog.Info("Starting in Local Development mode.")
		loadEnvVars(envFile)
	}

	if debugProfile {
		slog.Info("Starting application in debug mode,")
		EnableProfiling()
		setLoggingLevel(slog.LevelDebug)
		slog.Debug("Logging level set to Debug")
	}

	if dbUrl == "" {
		dbUrl = os.Getenv("DATABASE_URL")
	}

	dbConn := initPgConnPool(dbUrl)
	api.SetSwaggerSpec(swaggerSpec)
	healthCheckService := healthcheck.HealthCheckServicePgxImpl{DbConn: dbConn}
	auth_svc := s2auth.LocalAuthService{DbConn: dbConn}
	userCrud_svc := s2usercrud.UserCrudPgxImpl{DbConn: dbConn}
	if resetDevUSer {
		uid, err := uuid.Parse("b0fed113-30c4-42aa-bcdb-d0ecf2ac7f14")
		if err != nil {
			log.Fatalf("Error Parsing UUID for admin user %s\n", err.Error())
		}
		slog.Info("Reseting devuser Admin accoutn devuser pw")
		userCrud_svc.UpdateUserPasswordById(uid, os.Getenv("DEV_APP_PASS"))

	}
	err := api.StartApiServer(srvport, &healthCheckService, &auth_svc, &userCrud_svc)
	if err != nil {
		slog.Error("error creating new server instance", slog.String("error", err.Error()))
	}
}
