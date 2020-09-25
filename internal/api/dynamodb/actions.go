package db

import (
	"math/rand"

	"github.com/Travmatth-org/qaas/internal/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
    "github.com/aws/aws-sdk-go/service/dynamodb" 
)

func (d *DynamoDBClient) getPutItemInput(table *string, m interface{}, av *types.AWSAttrVal) (*types.AWSGetItemInput, error) {
	var cond types.AWSConditionBuilder
	switch m.(type) {
	case types.Quote:
		cond = expression.Name("ID").AttributeNotExists()
	default:
		notQuote := expression.Name("QuoteID").AttributeNotExists()
		notName := expression.Name("Name").AttributeNotExists()
		cond = expression.And(notQuote, notName)
	}
	expr, err := expression.NewBuilder().WithCondition(cond).Build()
	// TODO delete after verify
	if err != nil {
		return nil, err
	}
	return (&dynamodb.GetItemInput{}).
		SetItem(av).
		SetConditionExpression(expr.KeyCondition()).
		SetTableName(table), nil
}

func (d *DynamoDBClient) PutItem(m interface{}, table *string) error {
	av, err := dynamodbattribute.MarshalMap(m)
	// TODO delete after verify
	if err != nil {
		return err
	}
	_, err = d.client.PutItem(d.getPutItemInput(table, m, av))
	// TODO log
	return err
}

func (d *DynamoDBClient) getQuoteFromItems(items types.AWSAttrValMap) (*types.Quote, error) {
	id := (&dynamodb.GetItemInput{}).SetKey(items[rand.Intn(5)])
	r, err := d.client.GetItem(id)
	// TODO log
	if err != nil {
		return nil, err
	}
	var q types.Quote
	err = dynamodbattribute.UnmarshalMap(r.Item, &q)
	// TODO delete after verify
	if err != nil {
		return err
	}
	return &q, err
}

// GetRandomQuote returns a random quote from the database by scanning
// over the quotes table and projecting the entries into an array of [id].
// A random Id is selected from the list, and the quote fetched
func (d *DynamoDBClient) GetRandomQuote() *types.QuoteResponse {
	input := (&dynamodb.ScanInput{}).
		SetTableName(d.QuoteTable).
		SetProjectionExpression("ID").
		SetLimit(aws.Int(5))
	res, err := d.client.Scan(input())
	if err != nil {
		return &types.QuoteResponse{nil, err}
	}
	quote, err := d.getQuoteFromItems(res.Items)
	return types.NewQuoteResponse().WithQuote(quote).WithErr(err)
}

func (d *DynamoDBClient) buildQueryInput(table *string, expr expression.Expression, av *types.AWSAttrValMap) *types.AWSQueryInput {
	return (&dynamodb.QueryInput{}).
		SetLimit(d.limit).
		SetProjectionExpression(expr.Projection()).
		SetTableName(table).
		SetFilterExpression(expr.Filter()).
		SetExpressionAttributeNames(expr.Names()).
		SetExpressionAttributeValues(expr.Values()).
		SetReturnConsumedCapacity(aws.String("INDEXES")).
		SetExclusiveStartKey(av)
}

func (d *DynamoDBClient) getQueryInput(table, attr, start *string) (*types.AWSQueryInput, error) {
	filt := expression.Equal(expression.Name("Name"), expression.Value(attr))
	expr, err := expression.NewBuilder().
		WithFilter(filt).
		WithProjection(expression.NamesList(expression.Name("QuoteID"))).
		Build()
	// TODO delete after verify
	if err != nil {
		return nil, err
	}
	av, err := dynamodb.MarshalMap(struct{ ID string }{*start})
	// TODO delete after verify
	if err != nil {
		return nil, err
	}
	return d.buildQueryInput(table, expr, av), nil
}

func (d *DynamoDBClient) BatchGetQuotes(table string, ids types.AWSAttrValMap) ([]*types.Quote, error) {
	items := (&dynamodb.KeysAndAttributes{}).SetKeys(ids)
	input := (&dynamodb.BatchGetItemInput{}).
		SetRequestItems(items).
		SetReturnConsumedCapacity(true)
	res, err := d.client.BatchGetItem(input)
	// TODO log
	if err != nil {
		return []*types.Quote{}, err
	}
	var quotes []*types.Quote
	err = dynamodbattribute.UnmarshalList(res.Responses[d.QuoteTable], &quotes)
	return quotes, err
}

func (d *DynamoDBClient) GetQuotesByAttr(table, attr, start *string) *types.MultiQuoteResponse {
	input, err := d.getQueryInput(table, attr, start)
	// TODO delete after verify
	if err != nil {
		return types.NewMultiQuoteResponse().WithErr(err)
	}
	res, err := d.client.Query(input)
	// TODO log
	if err != nil {
		return types.NewMultiQuoteResponse().WithErr(err)
	}
	if quotes, err := d.client.BatchGetQuotes(table, res.Items); err != nil {
		return types.NewMultiQuoteResponse().WithErr(err)
	} else {
		last := res.LastEvaluatedKey["ID"]["S"]
		return types.NewMultiQuoteResponse().WithQuotes(quotes).WithNext(last)
	}
}
