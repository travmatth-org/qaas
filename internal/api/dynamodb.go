package api

import (
	db "github.com/travmatth-org/qaas/internal/api/dynamodb"
	"github.com/travmatth-org/qaas/internal/types"
	"honnef.co/go/tools/config"
)

// NewDynamoDBClient constructs a client for dynamodb communication
func (a *API) NewDynamoDBClient(c *config.Config) {
	a.dbClient = db.New().
		Init(a.session, a.config).
		SetPaginationLimit(c.PaginationLimit).
		SetQuoteTable(c.QuoteTable).
		SetTopicTable(c.TopicTable).
		SetAuthorTable(c.AuthorTable)
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

func (a *API) GetRandomQuote() *types.Quote {
	return a.dbClient.GetRandomQuote()
}

func (a *API) GetQuotesByTopic(topic, start *string) *types.MultiQuoteResponse {
	return a.dbClient.GetQuotesByAttr(a.dbClient.TopicTable, topic, start)
}

func (a *API) GetQuotesByAuthor(author, start *string) *types.MultiQuoteResponse {
	return a.dbClient.GetQuotesByAttr(a.dbClient.AuthorTable, author, start)
}
