package main

import (
	"os"

	"github.com/travmatth-org/qaas/internal/config"
	"github.com/travmatth-org/qaas/internal/logger"
	"github.com/travmatth-org/qaas/internal/server"

	"github.com/aws/aws-xray-sdk-go/awsplugins/ec2"
	"github.com/aws/aws-xray-sdk-go/xray"
)

func main() {
	c := config.Build()
	if c == nil {
		logger.Error().Msg("Error configuring server")
		os.Exit(1)
	}

	if c.IsProd() {
		ec2.Init()
		err := xray.Configure(xray.Config{ServiceVersion: "1.2.3"})
		if err != nil {
			logger.Error().Err(err).Msg("Error initializing x-ray tracing")
			os.Exit(1)
		}
	}

	// Create server
	s := server.New(c)
	s.RegisterHandlers()

	// run server
	if err := s.AcceptConnections(); err != nil {
		os.Exit(1)
	}
}
