package api

import (
	"errors"

	"github.com/travmatth-org/qaas/internal/logger"
	"github.com/travmatth-org/qaas/internal/types"
)

// PutNewQuote enters a new quote into DynamoDB
func (a *API) PutNewQuote(quote *types.Quote, author *types.Record, topics []*types.Record) error {
	if a.DDB == nil {
		err := errors.New("Error: DynamoDB agent not initialized")
		return err
	}
	var (
		txIn  = a.DDB.NewQuoteTxIn()
		table = a.DDB.Table
	)

	err := a.DDB.AddTxItem(txIn, quote, table.Quote)
	if err != nil {
		return err
	}
	err = a.DDB.AddTxItem(txIn, author, table.Author)
	if err != nil {
		logger.Error().Err(err).Msg("Error creating new transaction")
		return err
	}

	for _, topic := range topics {
		err := a.DDB.AddTxItem(txIn, topic, table.Topic)
		if err != nil {
			logger.Error().Err(err).Msg("Error creating new transaction")
			return err
		}
	}
	return a.DDB.PutNewQuoteTx(txIn)
}

// GetQuotesByAttr fetches a paginated response of quote related to attribute
func (a *API) GetQuotesByAttr(table, attr string, last *types.Record) *types.MultiQuoteRes {
	expr, err := a.DDB.NewQueryExpr(attr)
	if err != nil {
		logger.Error().Err(err).Msg("Error creating new query expression")
		return types.NewMultiQuoteRes().WithErr(err)
	}

	queryIn := a.DDB.NewQueryInput(expr, table)
	if last != nil {
		queryIn = a.DDB.SetQueryInStart(queryIn, last)
	}
	queryOut, err := a.DDB.QueryObject(queryIn)
	if err != nil {
		logger.Error().Err(err).Msg("Error querying for quotes")
		return types.NewMultiQuoteRes().WithErr(err)
	}

	items, last := a.DDB.ProcessQueryObjects(queryOut)
	if len(items) == 0 {
		return types.NewMultiQuoteRes()
	}
	batchIn := a.DDB.NewBatchGetInput(items)
	batchOut, err := a.DDB.BatchGetQuotes(batchIn)
	if err != nil {
		logger.Error().Err(err).Msg("Error requesting quotes")
		return types.NewMultiQuoteRes().WithErr(err)
	}

	quotes := a.DDB.ExtractQuotes(batchOut)
	if last != nil {
		return quotes.WithNext(last)
	}
	return quotes
}

// GetQuotesByTopic retrieves a paginated response of quote related to topic
func (a *API) GetQuotesByTopic(topic string, last *types.Record) *types.MultiQuoteRes {
	return a.GetQuotesByAttr(a.DDB.Table.Topic, topic, last)
}

// GetQuotesByAuthor retrieves a paginated response of quote related to author
func (a *API) GetQuotesByAuthor(author string, last *types.Record) *types.MultiQuoteRes {
	return a.GetQuotesByAttr(a.DDB.Table.Author, author, last)
}

// GetRandomQuote fetches a random quote from the database
func (a *API) GetRandomQuote() *types.QuoteRes {
	scanIn, err := a.DDB.NewScanIn()
	if err != nil {
		return types.NewQuoteRes().WithErr(err)
	}
	scanOut, err := a.DDB.ScanIDs(scanIn)
	if err != nil {
		return types.NewQuoteRes().WithErr(err)
	}

	id := a.DDB.RandomIDFromScan(scanOut)
	getIn := a.DDB.NewGetQuoteInputByID(id)
	getOut, err := a.DDB.GetObject(getIn)
	if err != nil {
		return types.NewQuoteRes().WithErr(err)
	}

	quote, err := a.DDB.QuoteFromObject(getOut)
	if err != nil {
		return types.NewQuoteRes().WithErr(err)
	}
	return types.NewQuoteRes().WithQuote(quote)
}
