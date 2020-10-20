package api

import (
	"errors"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/travmatth-org/qaas/internal/logger"
	"github.com/travmatth-org/qaas/internal/types"
)

// PutInput represents the input of the Put handler
type PutInput struct {
	input *types.DBTxWrIn
}

// PutOpt represents the configuration for PutInput
type PutOpt func(*PutInput) (*PutInput, error)

// WithItem adds an item to to the transaction
func (p *PutInput) WithItem(v interface{}, t string) (*PutInput, error) {
	av, err := dynamodbattribute.MarshalMap(v)
	if err != nil {
		logger.Error().Err(err).Msg("Error creating new transaction")
		return nil, err
	}
	in := &types.DBTxWrItem{
		Put: &dynamodb.Put{
			Item:      av,
			TableName: &t,
		}}
	p.input.TransactItems = append(p.input.TransactItems, in)
	return p, nil
}

// PutWithQuote adds quote to transaction
func (a *API) PutWithQuote(quote *types.Quote) PutOpt {
	return func(p *PutInput) (*PutInput, error) {
		return p.WithItem(quote, a.Table.Quote)
	}
}

// PutWithAuthor add author to transaction
func (a *API) PutWithAuthor(author *types.Record) PutOpt {
	return func(p *PutInput) (*PutInput, error) {
		return p.WithItem(author, a.Table.Quote)
	}
}

// PutWithTopics adds topics to the transaction
func (a *API) PutWithTopics(topics []*types.Record) PutOpt {
	return func(p *PutInput) (*PutInput, error) {
		for _, topic := range topics {
			switch in, err := p.WithItem(topic, a.Table.Topic); {
			case err != nil:
				logger.Error().Err(err).Msg("Error creating new transaction")
				return nil, err
			default:
				p = in
			}
		}
		return p, nil
	}
}

// helper function to iterate over and execute options
func buildPut(opts ...PutOpt) (*PutInput, error) {
	var (
		in = &PutInput{
			input: &types.DBTxWrIn{
				ReturnConsumedCapacity:      &AWSConsumedCapacityIndex,
				ReturnItemCollectionMetrics: &AWSItemCollection,
				TransactItems:               []*types.DBTxWrItem{},
			}}
		err error = nil
	)

	for _, opt := range opts {
		in, err = opt(in)
		if err != nil {
			return nil, err
		}
	}
	return in, err
}

// Put enters a new quote into DynamoDB
func (a *API) Put(opts ...PutOpt) error {
	if a.service == nil {
		return errors.New("Error: DynamoDB agent not initialized")
	}

	p, err := buildPut(opts...)
	if err != nil {
		return err
	}
	out, err := a.service.TransactWriteItems(p.input)
	logger.InfoOrErr(err).Interface("out", out).Msg("Submitted new quote entry")
	return err
}
