package db

import (
	"math/rand"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/travmatth-org/qaas/internal/types"
)

// getQuoteFromItems retreives a random quote item from a list of possible items
func (d *DynamoDBClient) getQuoteFromItems(items types.AWSAttrValMap) *types.QuoteResponse {
	id := (&dynamodb.GetItemInput{}).SetKey(items[rand.Intn(5)])
	// TODO log
	r, err := d.client.GetItem(id)
	if err != nil {
		return types.NewQuoteResponse().WithErr(err)
	}
	var q types.Quote
	err = dynamodbattribute.UnmarshalMap(r.Item, &q)
	return types.NewQuoteResponse().WithQuote(&q).WithErr(err)
}

// GetRandomQuote returns a random quote from the database by scanning
// over the quotes table and projecting the entries into an array of [id].
// A random Id is selected from the list, and the quote fetched
func (d *DynamoDBClient) GetRandomQuote() *types.QuoteResponse {
	res, err := d.client.Scan((&dynamodb.ScanInput{}).
		SetTableName(d.QuoteTable).
		SetProjectionExpression("ID").
		SetLimit(5))
	if err != nil {
		return types.NewQuoteResponse().WithErr(err)
	}
	return d.getQuoteFromItems(res.Items)
}

// getQueryInput configures a dynamodb getinput struct for a get request
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

// BatchGetQuotes retreives a reponse quote set specified by ID
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

// GetQuotesByAttr retrieves a paginated reponse of items from table
// by attribute, starting at start token
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
	quotes, err := d.BatchGetQuotes(table, res.Items)
	if err != nil {
		return types.NewMultiQuoteResponse().WithErr(err)
	}
	last := *res.LastEvaluatedKey["ID"].S
	return types.NewMultiQuoteResponse().WithQuotes(quotes).WithNext(last)
}

// GetQuotesByTopic retreives paginated response of quotes
// in given topic, starting at start token
func (d *DynamoDBClient) GetQuotesByTopic(topic, start string) *types.MultiQuoteResponse {
	return d.GetQuotesByAttr(d.TopicTable, topic, start)
}

// GetQuotesByAuthor retrieves paginated response of quotes
// by given author, starting at start token
func (d *DynamoDBClient) GetQuotesByAuthor(author, start string) *types.MultiQuoteResponse {
	return d.GetQuotesByAttr(d.AuthorTable, author, start)
}

// PutItem enters an item into the given table
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

// PutQuote enters a new quote, adding Quote, Author, and Topics to tables
func (d *DynamoDBClient) PutQuote(quote, from string, topics []string) error {
	q := types.NewQuote().WithText(quote).WithAuthor(from).WithTopics(topics)
	if err := d.PutItem(q, d.QuoteTable); err != nil {
		return err
	}
	for _, topic := range q.Topics {
		t := types.NewTopic().WithName(topic).WithQuoteID(q.ID)
		err := d.PutItem(t, d.TopicTable)
		if err != nil {
			return err
		}
	}
	author := types.NewAuthor().WithName(from).WithQuoteID(q.ID)
	return d.PutItem(author, d.AuthorTable)
}
