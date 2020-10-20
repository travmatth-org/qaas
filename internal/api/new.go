package api

import (
	"errors"

	"github.com/travmatth-org/qaas/internal/config"
	"github.com/travmatth-org/qaas/internal/logger"
	"github.com/travmatth-org/qaas/internal/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-xray-sdk-go/awsplugins/ec2"
	"github.com/aws/aws-xray-sdk-go/xray"
)

var (
	// DBRoleARN specifies the IAM Role to assume
	DBRoleARN = ""
	// AWSConsumedCapacityIndex specifies DynamoDB return aggregate RCU usage
	AWSConsumedCapacityIndex = dynamodb.ReturnConsumedCapacityIndexes
	// AWSItemCollection specifies the metric types returned
	AWSItemCollection = dynamodb.ReturnItemCollectionMetricsSize
)

// API manages client connections with outside services
type API struct {
	Table   config.Tables
	limit   int64
	service types.DBIFace
}

// Opts is the type signature for optional functions modifying API
type Opts func(*API) (*API, error)

// New constructs and returns an api client for client communications
func New(opts ...Opts) (*API, error) {
	var (
		err error = nil
		a         = &API{}
	)

	for _, opt := range opts {
		if a, err = opt(a); err != nil {
			return nil, err
		}
	}

	if a.service == nil {
		logger.Error().Msg("Error creating DB, nil agent")
		err = errors.New("Initialization failed, nil Service")
	}
	return a, err
}

// WithSession inserts a given session into API
func WithSession(a *API) (*API, error) {
	return a, nil
}

// WithPaginationLimit inserts the given pagination limit into DB
func WithPaginationLimit(limit int64) Opts {
	return func(a *API) (*API, error) {
		a.limit = limit
		return a, nil
	}
}

// WithTables inserts the name of the tables into DB
func WithTables(t config.Tables) Opts {
	return func(a *API) (*API, error) {
		a.Table = t
		return a, nil
	}
}

// WithEC2 inits EC2 client if production
func WithEC2(isProd bool) Opts {
	return func(a *API) (*API, error) {
		if isProd {
			ec2.Init()
		}
		return a, nil
	}
}

// WithXray inits XRay tracing if production
func WithXray(isProd bool) Opts {
	return func(a *API) (*API, error) {
		if !isProd {
			return a, nil
		}
		return a, xray.Configure(xray.Config{ServiceVersion: "1.2.3"})
	}
}

// WithDynamoDBService initializes the aws dynamodb agent
func WithDynamoDBService(c *config.Config) Opts {
	return func(a *API) (*API, error) {
		var (
			isProd   = config.IsProd(c)
			region   = c.AWS.Region
			endpoint = c.AWS.DynamoDB.Endpoint
			config   = aws.NewConfig().WithRegion(region)
		)

		sess, err := session.NewSession()
		switch {
		case isProd && DBRoleARN == "":
			err = errors.New("Error: DBRoleARN expected but empty")
			logger.Error().Err(err).Msg("Error initializing DB")
			fallthrough
		case err != nil:
			return nil, err
		case isProd:
			creds := stscreds.NewCredentials(sess, DBRoleARN)
			config = config.WithCredentials(creds)
		default:
			config = config.WithEndpoint(endpoint)
		}
		a.service = types.DBIFace(dynamodb.New(sess, config))
		return a, nil
	}
}
