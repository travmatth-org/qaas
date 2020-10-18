package types

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

// internal/api/api.go

// AWSSession ...
type AWSSession = session.Session

// AWSConfig ...
type AWSConfig = aws.Config

// AWSSTSCreds ...
type AWSSTSCreds = credentials.Credentials

// internal/api/dynamodb/actions

// AWSDynamoDB ...
type AWSDynamoDB = dynamodb.DynamoDB

// AWSQueryInput ...
type AWSQueryInput = dynamodb.QueryInput

// AWSAttrVal ...
type AWSAttrVal = dynamodb.AttributeValue

// AWSConditionBuilder ...
type AWSConditionBuilder = expression.ConditionBuilder

// AWSAttrValMap ...
type AWSAttrValMap = []map[string]*dynamodb.AttributeValue

// AWSGetItemInput ...
type AWSGetItemInput = dynamodb.GetItemInput

// AWSDDBIFace ...
type AWSDDBIFace = dynamodbiface.DynamoDBAPI

// AWSDDBTxWrItem ...
type AWSDDBTxWrItem = dynamodb.TransactWriteItem

// AWSDDBTxWrIn ...
type AWSDDBTxWrIn = dynamodb.TransactWriteItemsInput

// AWSDDBTxWrOut ...
type AWSDDBTxWrOut = dynamodb.TransactWriteItemsOutput

// AWSDDBGetOut ...
type AWSDDBGetOut = dynamodb.GetItemOutput

// AWSDDBGetIn ...
type AWSDDBGetIn = dynamodb.GetItemInput

// AWSDDBBatchGetIn ...
type AWSDDBBatchGetIn = dynamodb.BatchGetItemInput

// AWSDDBBatchGetOut ...
type AWSDDBBatchGetOut = dynamodb.BatchGetItemOutput

// AWSDDBAV ...
type AWSDDBAV = dynamodb.AttributeValue

// AWSDDBQueryIn ...
type AWSDDBQueryIn = dynamodb.QueryInput

// AWSDDBQueryOut ...
type AWSDDBQueryOut = dynamodb.QueryOutput

// AWSDDBExpr ...
type AWSDDBExpr = expression.Expression
