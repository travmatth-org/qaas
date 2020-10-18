package api

import (
	"github.com/travmatth-org/qaas/internal/api/ddb"
	"github.com/travmatth-org/qaas/internal/config"
	"github.com/travmatth-org/qaas/internal/types"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-xray-sdk-go/awsplugins/ec2"
	"github.com/aws/aws-xray-sdk-go/xray"
)

// API manages client connections with outside services
type API struct {
	session *types.AWSSession
	region  string
	DDB     *ddb.DDB
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
	return a, err
}

// WithRegion inserts a given region into API
func WithRegion(r string) Opts {
	return func(a *API) (*API, error) {
		a.region = r
		return a, nil
	}
}

// WithSession inserts a given session into API
func WithSession(a *API) (*API, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	a.session = sess
	return a, nil
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

// WithNewDDB configures and inserts DDB into API
func WithNewDDB(c *config.Config) Opts {
	return func(a *API) (*API, error) {
		var (
			isProd  = config.IsProd(c)
			db, err = ddb.New(
				ddb.WithAWSConfig(a.region),
				ddb.WithConfigEndpoint(c.AWS.DynamoDB.Endpoint, isProd),
				ddb.WithSTSCreds(isProd),
				ddb.WithAWSSession(a.session),
				ddb.WithPaginationLimit(c.AWS.DynamoDB.PaginationLimit),
				ddb.WithTables(c.AWS.DynamoDB.Table),
				ddb.NewClient,
			)
		)
		if err != nil {
			return nil, err
		}
		a.DDB = db
		return a, nil
	}
}

// WithDDB inserts a client into the API, useful for testing
func WithDDB(d *ddb.DDB) Opts {
	return func(a *API) (*API, error) {
		a.DDB = d
		return a, nil
	}
}
