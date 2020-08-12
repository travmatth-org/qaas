package main

import (
	"os"

	"github.com/Travmatth/faas/internal/config"
	"github.com/Travmatth/faas/internal/logger"
	"github.com/Travmatth/faas/internal/server"

	"github.com/aws/aws-xray-sdk-go/awsplugins/ec2"
	"github.com/aws/aws-xray-sdk-go/xray"
)

func init() {
	// conditionally load plugin
	if os.Getenv("ENVIRONMENT") == "production" {
		ec2.Init()
	}

	err := xray.Configure(xray.Config{ServiceVersion: "1.2.3"})
	if err != nil {
		logger.Error().Err(err).Msg("Error initializing x-ray tracing")
		os.Exit(1)
	}
}

func main() {
	// Config vals of server
	c := config.Build()
	if c == nil {
		logger.Error().Msg("Error configuring server")
		os.Exit(1)
	}

	// Create server
	s := server.New(c)
	s.RegisterHandlers()

	// run server
	os.Exit(s.AcceptConnections())
}
