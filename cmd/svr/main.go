// Package classification Golang Web App Example
//
// Golang Web App Example
//
//     Schemes: https
//     BasePath: /
//     Version: 1.0.0
//     Contact: Dustin Ratcliffe<dustin.k.ratcliffe@gmail.com>
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - Bearer:
//
//     SecurityDefinitions:
//     Bearer:
//          type: apiKey
//          name: Authorization
//          in: header
// swagger:meta
package main

import (
	"time"

	"github.com/DustyRat/go-webapp/internal/server"
	"github.com/DustyRat/go-webapp/internal/service"

	_ "github.com/go-openapi/errors"   // Go Open API errors
	_ "github.com/go-openapi/strfmt"   // Go Open API fmt
	_ "github.com/go-openapi/swag"     // Go Open API swag
	_ "github.com/go-openapi/validate" // Go Open API validate
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

var (
	buildDate,
	buildHost,
	gitURL,
	branch,
	sha string
	start = time.Now()
)

func main() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	err := server.Run(&service.BuildInfo{
		BuildDate: buildDate,
		BuildHost: buildHost,
		GitURL:    gitURL,
		Branch:    branch,
		SHA:       sha,
		Start:     start,
	})
	if err != nil {
		log.Fatal().Stack().Caller().Err(err).Send()
	}
}
