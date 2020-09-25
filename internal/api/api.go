package api

import (
	"github.com/aws/aws-sdk-go/aws/session"
	db "github.com/travmatth-org/qaas/internal/api/dynamodb"
	"github.com/travmatth-org/qaas/internal/config"
	"github.com/travmatth-org/qaas/internal/types"
)

// API manages client connections with outside services
type API struct {
	session  *types.AWSSession
	region   string
	env      string
	dbClient *db.DynamoDBClient
}

// New constructs and returns an api client for client communications
func New(c *config.Config) *API {
	return &API{}
}

func (a *API) WithSession() *API {
	a.session = session.Must(session.NewSession())
	return a
}

func (a *API) WithRegion(r string) *API {
	a.region = r
	return a
}

func (a *API) WithEnv(e string) *API {
	a.env = e
	return a
}
