package api

import (
	db "github.com/travmatth-org/qaas/internal/api/dynamodb"
	"github.com/travmatth-org/qaas/internal/config"
	"github.com/travmatth-org/qaas/internal/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/aws/credentials/stscreds" 
)

// API manages client connections with outside services
type API struct {
	session  *types.AWSSession
	region   string
	env      string
	config   *types.AWSConfig
	dbClient *db.DynamoDBClient
}

// New constructs and returns an api client for client communications
func New(c *config.Config) *API {
	return &API{}
}

func (a *API) WithRegion(r string) *API {
	a.region = r
	return a
}

func (a *API) WithEnv(e string) *API {
	a.env = e
	return a
}

func (a *API) WithCreds(arn, endpoint string) *API {
	a.session = session.Must(session.NewSession())
	a.config = aws.NewConfig().WithRegion(a.region)
	if a.env == "PRODUCTION" {
		a.config = a.config.WithEndpoint(endpoint)
	} else {
		creds := stscreds.NewCredentials(a.session, arn)
		a.config = a.config.WithCredentials(creds)
	}
	return a
}
