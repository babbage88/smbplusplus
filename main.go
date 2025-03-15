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
	"log/slog"

	"github.com/babbage88/smbplusplus/database/s2_pgxpool"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed swagger.yaml
var swaggerSpec []byte

func initPgConnPool() *pgxpool.Pool {
	connPool := s2_pgxpool.PgPoolInit()
	return connPool
}

func main() {
	//dbConn := initPgConnPool()
	server, err := NewSmbPlusServerFromConfig(".env")
	if err != nil {
		slog.Error("error creating new server instance", slog.String("error", err.Error()))

	}
	server.Start()
}
