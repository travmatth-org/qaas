package api

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/travmatth-org/qaas/internal/logger"
	"github.com/travmatth-org/qaas/internal/types"
)

// construct and perform DynamoDB query on given table, project only ID attr
func (a *API) query(table, name string, last *types.Record) (*types.DBQueryOut, error) {
	var (
		expr, _ = expression.NewBuilder().
			WithKeyCondition(expression.Key("Name").Equal(expression.Value(name))).
			WithProjection(expression.NamesList(expression.Name("QuoteID"))).
			Build()
		in = &types.DBQueryIn{
			ProjectionExpression:      expr.Projection(),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			Limit:                     &a.limit,
			ReturnConsumedCapacity:    &AWSConsumedCapacityIndex,
			TableName:                 &table,
			KeyConditionExpression:    expr.KeyCondition(),
		}
	)

	if last != nil {
		m, _ := dynamodbattribute.MarshalMap(last)
		in = in.SetExclusiveStartKey(m)
	}
	return a.service.Query(in)
}

// Extract quote objects from query response and map to quote ID AttributeValue map
func queryToIDs(out *types.DBQueryOut) ([]map[string]*types.DBAV, *types.Record) {
	r := types.NewRecord()
	for _, item := range out.Items {
		item["ID"] = item["QuoteID"]
		delete(item, "QuoteID")
	}
	if out.LastEvaluatedKey != nil {
		err := dynamodbattribute.UnmarshalMap(out.LastEvaluatedKey, r)
		if err != nil {
			logger.Error().Err(err).Msg("Error unmarshalling last key, setting nil")
		}
	}
	return out.Items, r
}

// configure and perform BatchGetItem operation
func (a *API) batch(avs []map[string]*types.DBAV) (*types.DBBatchGetOut, error) {
	var (
		in = &types.DBBatchGetIn{
			RequestItems: map[string]*dynamodb.KeysAndAttributes{
				a.Table.Quote: {
					Keys: avs,
				},
			},
			ReturnConsumedCapacity: &AWSConsumedCapacityIndex,
		}
	)
	return a.service.BatchGetItem(in)
}

// Get fetches a paginated response of quotes related to attribute
func (a *API) Get(table, attr string, last *types.Record) *types.MultiQuoteRes {
	var (
		res    = types.NewMultiQuoteRes()
		quotes = []*types.Quote{}
		t      = a.Table.Quote
	)

	in, err := a.query(table, attr, last)
	logger.InfoOrErr(err).Msg("Queried for quotes")
	if err != nil {
		return res.WithErr(err)
	}

	items, last := queryToIDs(in)
	if len(items) == 0 {
		return res.WithErr(fmt.Errorf("No Quotes Matching: %s", attr))
	}

	out, err := a.batch(items)
	logger.InfoOrErr(err).Interface("out", out).Msg("Error querying for quotes")
	if err != nil {
		return res.WithErr(err)
	}

	err = dynamodbattribute.UnmarshalListOfMaps(out.Responses[t], &quotes)
	return res.WithQuotes(quotes).WithErr(err).WithNext(last)
}
