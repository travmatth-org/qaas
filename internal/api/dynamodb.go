package api

import (
	db "github.com/travmatth-org/qaas/internal/api/dynamodb"
	"github.com/travmatth-org/qaas/internal/config"
	"github.com/travmatth-org/qaas/internal/types"
)

// DBRoleARN injected during build process
var DBRoleARN string

// NewDynamoDBClient constructs a client for dynamodb communication
func (a *API) NewDynamoDBClient(c *config.Config) {
	isProd := c.Env == config.Production
	a.dbClient = db.New(
		db.WithAWSConfig(a.region),
		db.WithConfigEndpoint(c.AWS.DynamoDB.Endpoint, isProd),
		db.WithSTSCreds(DBRoleARN, isProd),
		db.WithAWSSession(a.session),
		db.WithPaginationLimit(c.AWS.DynamoDB.PaginationLimit),
		db.WithQuoteTable(c.AWS.DynamoDB.Table.Quote),
		db.WithTopicTable(c.AWS.DynamoDB.Table.Topic),
		db.WithAuthorTable(c.AWS.DynamoDB.Table.Author),
	)
}

// PutQuote enters a new quote, adding Quote, Author, and Topics to tables
func (a *API) PutQuote(quote, from string, topics []string) error {
	q := types.NewQuote(WithText(quote), WithAuthor(from), WithTopics(topics))
	if err := a.dbClient.PutItem(q, a.dbClient.QuoteTable); err != nil {
		return err
	}
	for _, topic := range q.Topics {
		t := types.NewTopic().WithName(topic).WithQuoteID(q.ID)
		err = a.dbClient.PutItem(t, a.dbClient.TopicTable)
		if err != nil {
			return err
		}
	}
	author := types.NewAuthor().WithName(from).WithQuoteID(q.ID)
	return a.dbClient.PutItem(author, a.dbClient.AuthorTable)
}

func (a *API) GetRandomQuote() *types.QuoteResponse {
	return a.dbClient.GetRandomQuote()
}

func (a *API) GetQuotesByTopic(topic, start string) *types.MultiQuoteResponse {
	return a.dbClient.GetQuotesByAttr(a.dbClient.TopicTable, topic, start)
}

func (a *API) GetQuotesByAuthor(author, start string) *types.MultiQuoteResponse {
	return a.dbClient.GetQuotesByAttr(a.dbClient.AuthorTable, author, start)
}
