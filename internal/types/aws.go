package types

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/aws/session" 
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
    "github.com/aws/aws-sdk-go/aws/credentials" 
)

// internal/api/api.go
type AWSSession = session.Session
type AWSConfig = aws.Config
type AWSSTSCreds = credentials.Credentials

// internal/api/dynamodb/actions
type AWSDynamoDB = dynamodb.DynamoDB
type AWSQueryInput = dynamodb.QueryInput
type AWSAttrVal = dynamodb.AttributeValue
type AWSConditionBuilder = expression.ConditionBuilder
type AWSAttrValMap = []map[string]*dynamodb.AttributeValue
type AWSGetItemInput = dynamodb.GetItemInput
