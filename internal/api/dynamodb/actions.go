package db

import (
	"math/rand"

	"github.com/travmatth-org/qaas/internal/types"
	// "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

// func (d *DynamoDBClient) getPutItemInput(table *string, av types.AWSAttrVal, cond types.AWSConditionBuilder) (*types.AWSGetItemInput, error) {
// }

func (d *DynamoDBClient) PutItem(m interface{}, table string) error {
	av, err := dynamodbattribute.MarshalMap(m)
	// TODO delete after verify
	if err != nil {
		return err
	}
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
		return err
	}
	input := (&dynamodb.PutItemInput{}).
		SetItem(av).
		SetConditionExpression(*expr.Condition()).
		SetTableName(table)
	_, err = d.client.PutItem(input)
	// TODO log
	return err
}

func (d *DynamoDBClient) getQuoteFromItems(items types.AWSAttrValMap) *types.QuoteResponse {
	var q types.Quote
	id := (&dynamodb.GetItemInput{}).SetKey(items[rand.Intn(5)])
	// TODO log
	r, err := d.client.GetItem(id)
	if err != nil {
		return types.NewQuoteResponse().WithErr(err)
	}
	err = dynamodbattribute.UnmarshalMap(r.Item, &q)
	return types.NewQuoteResponse().WithQuote(&q).WithErr(err)
}

// GetRandomQuote returns a random quote from the database by scanning
// over the quotes table and projecting the entries into an array of [id].
// A random Id is selected from the list, and the quote fetched
func (d *DynamoDBClient) GetRandomQuote() *types.QuoteResponse {
	input := (&dynamodb.ScanInput{}).
		SetTableName(d.QuoteTable).
		SetProjectionExpression("ID").
		SetLimit(5)
	if res, err := d.client.Scan(input); err != nil {
		return types.NewQuoteResponse().WithErr(err)
	} else {
		return d.getQuoteFromItems(res.Items)
	}
}

func (d *DynamoDBClient) getQueryInput(table, attr, start string) (*types.AWSQueryInput, error) {
	filt := expression.Equal(expression.Name("Name"), expression.Value(attr))
	expr, err := expression.NewBuilder().
		WithFilter(filt).
		WithProjection(expression.NamesList(expression.Name("QuoteID"))).
		Build()
	// TODO delete after verify
	if err != nil {
		return nil, err
	}
	av, err := dynamodbattribute.MarshalMap(struct{ ID string }{start})
	// TODO delete after verify
	if err != nil {
		return nil, err
	}
	return (&dynamodb.QueryInput{}).
		SetLimit(d.limit).
		SetProjectionExpression(*expr.Projection()).
		SetTableName(table).
		SetFilterExpression(*expr.Filter()).
		SetExpressionAttributeNames(expr.Names()).
		SetExpressionAttributeValues(expr.Values()).
		SetReturnConsumedCapacity("INDEXES").
		SetExclusiveStartKey(av), nil
}

func (d *DynamoDBClient) BatchGetQuotes(table string, ids types.AWSAttrValMap) ([]*types.Quote, error) {
	quotes := make([]*types.Quote, int(d.limit))
	requestedItems := make(map[string]*dynamodb.KeysAndAttributes)
	requestedItems[table] = (&dynamodb.KeysAndAttributes{}).SetKeys(ids)
	res, err := d.client.BatchGetItem((&dynamodb.BatchGetItemInput{}).
		SetRequestItems(requestedItems).
		SetReturnConsumedCapacity("INDEXES"))
	// TODO log
	if err != nil {
		return quotes, err
	}
	list := res.Responses[d.QuoteTable]
	err = dynamodbattribute.UnmarshalListOfMaps(list, &quotes)
	return quotes, err
}

func (d *DynamoDBClient) GetQuotesByAttr(table, attr, start string) *types.MultiQuoteResponse {
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
	if quotes, err := d.BatchGetQuotes(table, res.Items); err != nil {
		return types.NewMultiQuoteResponse().WithErr(err)
	} else {
		last := *res.LastEvaluatedKey["ID"].S
		return types.NewMultiQuoteResponse().WithQuotes(quotes).WithNext(last)
	}
}
