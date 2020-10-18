package ddb

import (
	"math"
	"math/rand"

	"github.com/travmatth-org/qaas/internal/logger"
	"github.com/travmatth-org/qaas/internal/types"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

// NewQuoteTxIn creates a new, empty quote transaction input
func (d *DDB) NewQuoteTxIn() *types.AWSDDBTxWrIn {
	return &types.AWSDDBTxWrIn{
		ReturnConsumedCapacity:      &awsConsumedCapacityIndex,
		ReturnItemCollectionMetrics: &awsItemCollection,
		TransactItems:               []*types.AWSDDBTxWrItem{},
	}
}

// AddTxItem adds given transaction item to the input
func (d *DDB) AddTxItem(in *types.AWSDDBTxWrIn, v interface{}, t string) error {
	av, err := dynamodbattribute.MarshalMap(v)
	if err != nil {
		logger.Error().Err(err).Msg("Error creating new transaction")
		return err
	}
	in.TransactItems = append(in.TransactItems, &types.AWSDDBTxWrItem{
		Put: &dynamodb.Put{
			Item:      av,
			TableName: &t,
		}})
	return nil
}

// PutNewQuoteTx places new into DynamoDB
func (d *DDB) PutNewQuoteTx(in *types.AWSDDBTxWrIn) error {
	logger.Info().Interface("in", in).Msg("Creating new quote entry")
	req, out := d.service.TransactWriteItemsRequest(in)
	err := req.Send()
	logger.InfoOrErr(err).Interface("out", out).Msg("Submitted new quote entry")
	return err
}

// NewGetQuoteInputByID Creates get action input for a given quote id
func (d *DDB) NewGetQuoteInputByID(id string) *types.AWSDDBGetIn {
	m, _ := dynamodbattribute.MarshalMap(map[string]string{"ID": id})
	return &types.AWSDDBGetIn{
		TableName: &d.Table.Quote,
		Key:       m,
	}
}

// GetObject retrieves given object from dynamodb with given input
func (d *DDB) GetObject(in *types.AWSDDBGetIn) (*types.AWSDDBGetOut, error) {
	logger.Info().Interface("in", in).Msg("Getting New Object")
	req, out := d.service.GetItemRequest(in)
	err := req.Send()
	logger.InfoOrErr(err).Interface("out", out).Msg("Completed Get Object")
	return out, err
}

// QuoteFromObject decodes the DynamoDB response into *Quote struct
func (d *DDB) QuoteFromObject(out *types.AWSDDBGetOut) (*types.Quote, error) {
	var (
		quote = types.NewQuote()
		err   = dynamodbattribute.UnmarshalMap(out.Item, quote)
	)
	return quote, err
}

// NewQueryExpr create *Expression for query input
func (d *DDB) NewQueryExpr(name string) (types.AWSDDBExpr, error) {
	return expression.NewBuilder().
		WithKeyCondition(expression.Key("Name").Equal(expression.Value(name))).
		WithProjection(expression.NamesList(expression.Name("QuoteID"))).
		Build()
}

// NewQueryInput creates a new QueryInput with the expression for the given table
func (d *DDB) NewQueryInput(expr types.AWSDDBExpr, table string) *types.AWSDDBQueryIn {
	return &types.AWSDDBQueryIn{
		ProjectionExpression:      expr.Projection(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Limit:                     &d.limit,
		ReturnConsumedCapacity:    &awsConsumedCapacityIndex,
		TableName:                 &table,
		KeyConditionExpression:    expr.KeyCondition(),
	}
}

// SetQueryInStart sets the start on for the given query
func (d *DDB) SetQueryInStart(in *types.AWSDDBQueryIn, last *types.Record) *types.AWSDDBQueryIn {
	m, err := dynamodbattribute.MarshalMap(last)
	if err != nil {
		logger.Error().Err(err).Msg("Unable to marshal last key, skipping")
		return in
	}
	return in.SetExclusiveStartKey(m)
}

// QueryObject queries DynamoDB for objects matching input
func (d *DDB) QueryObject(in *types.AWSDDBQueryIn) (*types.AWSDDBQueryOut, error) {
	logger.Info().Interface("in", in).Msg("Querying for Object")
	req, out := d.service.QueryRequest(in)
	err := req.Send()
	logger.InfoOrErr(err).Interface("out", out).Msg("Queried Object")
	return out, err
}

// ProcessQueryObjects processes DynamoDB query response into an array of attribute value maps
func (d *DDB) ProcessQueryObjects(out *types.AWSDDBQueryOut) ([]map[string]*types.AWSDDBAV, *types.Record) {
	var (
		r *types.Record = nil
	)
	for _, item := range out.Items {
		item["ID"] = item["QuoteID"]
		delete(item, "QuoteID")
	}
	if out.LastEvaluatedKey != nil {
		err := dynamodbattribute.UnmarshalMap(out.LastEvaluatedKey, &r)
		if err != nil {
			logger.Error().Err(err).Msg("Error unmarshalling last key, setting nil")
		}
	}
	return out.Items, r
}

// NewBatchGetInput creates new batch get input from the given items
func (d *DDB) NewBatchGetInput(avs []map[string]*types.AWSDDBAV) *types.AWSDDBBatchGetIn {
	return &types.AWSDDBBatchGetIn{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			d.Table.Quote: {
				Keys: avs,
			},
		},
		ReturnConsumedCapacity: &awsConsumedCapacityIndex,
	}
}

// BatchGetQuotes fetches quotes from DynamoDB
func (d *DDB) BatchGetQuotes(in *types.AWSDDBBatchGetIn) (*types.AWSDDBBatchGetOut, error) {
	logger.Info().Interface("in", in).Msg("Batch requesting quotes")
	req, out := d.service.BatchGetItemRequest(in)
	err := req.Send()
	logger.InfoOrErr(err).Interface("out", out).Msg("Requested Quote Batch")
	return out, err
}

// ExtractQuotes decodes DynamoDB response into a *MultiQuoteRes
func (d *DDB) ExtractQuotes(out *types.AWSDDBBatchGetOut) *types.MultiQuoteRes {
	var (
		quotes    = []*types.Quote{}
		responses = out.Responses[d.Table.Quote]
		err       = dynamodbattribute.UnmarshalListOfMaps(responses, &quotes)
	)
	return types.NewMultiQuoteRes().WithQuotes(quotes).WithErr(err)
}

// NewScanIn creates scan request input
func (d *DDB) NewScanIn() (*dynamodb.ScanInput, error) {
	proj := expression.NamesList(expression.Name("ID"))
	expr, err := expression.NewBuilder().WithProjection(proj).Build()
	if err != nil {
		logger.Error().Err(err).Msg("Error creating scan expression")
		return nil, err
	}
	return &dynamodb.ScanInput{
		TableName:                 &d.Table.Quote,
		ProjectionExpression:      expr.Projection(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}, nil
}

// ScanIDs scans dynamodb to for quotes IDs
func (d *DDB) ScanIDs(in *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
	logger.Info().Interface("in", in).Msg("Scanning for Quote IDs")
	req, out := d.service.ScanRequest(in)
	err := req.Send()
	logger.InfoOrErr(err).Interface("out", out).Msg("Requested Quote Batch")
	return out, err
}

// RandomIDFromScan returns a random ID from a Scans output
func (d *DDB) RandomIDFromScan(out *dynamodb.ScanOutput) string {
	ceiling := int(math.Min(float64(*out.Count), float64(d.limit)))
	return *out.Items[rand.Intn(ceiling)]["ID"].S
}
