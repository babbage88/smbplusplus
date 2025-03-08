// Package main go-infra API.
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

import _ "embed"

//go:embed swagger.yaml
var swaggerSpec []byte

func main() {
	server := NewSmbPlusServerFromConfig(".env")
	server.Start()
}
