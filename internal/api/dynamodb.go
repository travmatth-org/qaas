package api

import (
	db "github.com/travmatth-org/qaas/internal/api/dynamodb"
	"github.com/travmatth-org/qaas/internal/config"
	"github.com/travmatth-org/qaas/internal/types"
)

// NewDynamoDBClient constructs a client for dynamodb communication
func (a *API) NewDynamoDBClient(c *config.Config) {
	a.dbClient = db.New().WithAWSConfig(a.region).WithAWSSession(a.session)
	if c.IsProd() {
		a.dbClient = a.dbClient.WithSTSCreds(c.GetDBRoleARN())
	} else {
		a.dbClient = a.dbClient.WithConfigEndpoint(c.GetDBEndpoint())
	}
	a.dbClient = a.dbClient.
		WithPaginationLimit(c.PaginationLimit).
		WithQuoteTable(c.QuoteTable).
		WithTopicTable(c.TopicTable).
		WithAuthorTable(c.AuthorTable).
		NewClient()
}

// PutQuote enters a new quote, adding Quote, Author, and Topics to tables
func (a *API) PutQuote(quote, from string, topics []string) error {
	q := types.NewQuote().WithText(quote).WithAuthor(from).WithTopics(topics)
	err := a.dbClient.PutItem(q, a.dbClient.QuoteTable)
	if err != nil {
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
