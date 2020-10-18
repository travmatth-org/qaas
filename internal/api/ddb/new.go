package ddb

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/travmatth-org/qaas/internal/config"
	"github.com/travmatth-org/qaas/internal/logger"
	"github.com/travmatth-org/qaas/internal/types"
)

var (
	// DDBRoleARN specifies the IAM Role that the httpd server
	// assumes toperform DynamoDB service operations
	DDBRoleARN = ""
	// Specify that requests to DynamoDB service should
	// return info on capacity consumed in aggregate
	awsConsumedCapacityIndex = dynamodb.ReturnConsumedCapacityIndexes
	awsItemCollection        = dynamodb.ReturnItemCollectionMetricsSize
)

// DDB manages the communications with the DynamoDB service
type DDB struct {
	session *types.AWSSession
	config  *types.AWSConfig
	service types.AWSDDBIFace
	limit   int64
	// expr    *expression.Expression
	Table config.Tables
}

// Opts is the type signature of functions modifying a DDB struct
type Opts func(d *DDB) (*DDB, error)

// New creates and configures a new DDB with the specified options
func New(opts ...Opts) (*DDB, error) {
	var (
		d         = &DDB{}
		err error = nil
	)
	for _, opt := range opts {
		if d, err = opt(d); err != nil {
			return nil, err
		}
	}
	if d.service == nil {
		err = errors.New("Initialization failed, nil service")
		logger.Error().Err(err).Msg("Error creating DDB")
		return nil, err
	}
	return d, err
}

// WithAWSConfig inserts the given region in the aws config
func WithAWSConfig(region string) Opts {
	return func(d *DDB) (*DDB, error) {
		d.config = aws.NewConfig().WithRegion(region)
		return d, nil
	}
}

// WithAWSSession inserts the given AWS session object into DDB
func WithAWSSession(sess *types.AWSSession) Opts {
	return func(d *DDB) (*DDB, error) {
		d.session = sess
		return d, nil
	}
}

// WithSTSCreds configures the AWS config with STS creds using
// the injected role ARN, if running in production
func WithSTSCreds(isProd bool) Opts {
	return func(d *DDB) (*DDB, error) {
		if isProd && DDBRoleARN == "" {
			err := errors.New("Error: DDBRoleARN expected but empty")
			logger.Error().Err(err).Msg("Error initializing DDB")
			return nil, err
		} else if isProd {
			creds := stscreds.NewCredentials(d.session, DDBRoleARN)
			d.config = d.config.WithCredentials(creds)
		}
		return d, nil
	}
}

// WithConfigEndpoint configures the DynamoDB endpoint, if not running in prod
func WithConfigEndpoint(endpoint string, isProd bool) Opts {
	return func(d *DDB) (*DDB, error) {
		if !isProd {
			d.config = d.config.WithEndpoint(endpoint)
		}
		return d, nil
	}
}

// NewClient creates a new DynamoDB service using the provided session and config
func NewClient(d *DDB) (*DDB, error) {
	d.service = types.AWSDDBIFace(dynamodb.New(d.session, d.config))
	return d, nil
}

// WithClient inserts a dynamodb service into DDB, useful in testing
func WithClient(iface types.AWSDDBIFace) Opts {
	return func(d *DDB) (*DDB, error) {
		d.service = iface
		return d, nil
	}
}

// WithPaginationLimit inserts the given pagination limit into DDB
func WithPaginationLimit(limit int64) Opts {
	return func(d *DDB) (*DDB, error) {
		d.limit = limit
		return d, nil
	}
}

// WithTables inserts the name of the tables into DDB
func WithTables(t config.Tables) Opts {
	return func(d *DDB) (*DDB, error) {
		d.Table = t
		return d, nil
	}
}
