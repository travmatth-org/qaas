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

// DBIFace ...
type DBIFace = dynamodbiface.DynamoDBAPI

// DBTxWrItem ...
type DBTxWrItem = dynamodb.TransactWriteItem

// DBTxWrIn ...
type DBTxWrIn = dynamodb.TransactWriteItemsInput

// DBTxWrOut ...
type DBTxWrOut = dynamodb.TransactWriteItemsOutput

// DBGetOut ...
type DBGetOut = dynamodb.GetItemOutput

// DBGetIn ...
type DBGetIn = dynamodb.GetItemInput

// DBBatchGetIn ...
type DBBatchGetIn = dynamodb.BatchGetItemInput

// DBBatchGetOut ...
type DBBatchGetOut = dynamodb.BatchGetItemOutput

// DBAV ...
type DBAV = dynamodb.AttributeValue

// DBQueryIn ...
type DBQueryIn = dynamodb.QueryInput

// DBQueryOut ...
type DBQueryOut = dynamodb.QueryOutput

// DBExpr ...
type DBExpr = expression.Expression
