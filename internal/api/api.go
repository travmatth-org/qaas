package api

import (
	db "github.com/travmatth-org/qaas/internal/api/dynamodb"
	"github.com/travmatth-org/qaas/internal/config"
	"github.com/travmatth-org/qaas/internal/types"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-xray-sdk-go/awsplugins/ec2"
	"github.com/aws/aws-xray-sdk-go/xray"

)

// API manages client connections with outside services
type API struct {
	session  *types.AWSSession
	region   string
	DynamoDB *db.DynamoDBClient
}

type apiOpt func(*API) (*API, error)

// New constructs and returns an api client for client communications
func New(opts ...apiOpt) (*API, error) {
	var err error
	a := &API{}
	for _, opt := range opts {
		if a, err = opt(a); err != nil {
			return nil, err
		}
	}
	return a, nil
}

func WithRegion(r string) apiOpt {
	return func(a *API) (*API, error) {
		a.region = r
		return a, nil
	}
}

func WithSession(a *API) (*API, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	a.session = sess
	return a, nil
}

func WithEC2(isProd bool) apiOpt {
	return func(a *API) (*API, error) {
		if isProd {
			ec2.Init()
		}
		return a, nil
	}
}

func WithXray(isProd bool) apiOpt {
	return func(a *API) (*API, error) {
		if isProd {
			return a, nil
		}
		return a, xray.Configure(xray.Config{ServiceVersion: "1.2.3"})
	}
}

func WithNewDynamoDBClient(c *config.Config) apiOpt {
	return func(a *API) (*API, error) {
		isProd := c.Env == config.Production
		a.DynamoDB = db.New(
			db.WithAWSConfig(a.region),
			db.WithConfigEndpoint(c.AWS.DynamoDB.Endpoint, isProd),
			db.WithSTSCreds(isProd),
			db.WithAWSSession(a.session),
			db.WithPaginationLimit(c.AWS.DynamoDB.PaginationLimit),
			db.WithQuoteTable(c.AWS.DynamoDB.Table.Quote),
			db.WithTopicTable(c.AWS.DynamoDB.Table.Topic),
			db.WithAuthorTable(c.AWS.DynamoDB.Table.Author))
		return a, nil
	}
}

func WithDynamoDBClient(d *db.DynamoDBClient) apiOpt {
	return func(a *API) (*API, error) {
		a.DynamoDB = d
		return a, nil
	}
}