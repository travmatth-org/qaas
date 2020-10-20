package api

import (
	"errors"
	"math"
	"math/rand"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/travmatth-org/qaas/internal/logger"
	"github.com/travmatth-org/qaas/internal/types"
)

func (a *API) get(id map[string]*dynamodb.AttributeValue) (*types.DBGetOut, error) {
	var (
		in = &types.DBGetIn{
			TableName: &a.Table.Quote,
			Key:       id,
		}
	)
	return a.service.GetItem(in)
}

func (a *API) scan() (*dynamodb.ScanOutput, error) {
	var (
		proj    = expression.NamesList(expression.Name("ID"))
		expr, _ = expression.NewBuilder().WithProjection(proj).Build()
		in      = &dynamodb.ScanInput{
			TableName:                 &a.Table.Quote,
			ProjectionExpression:      expr.Projection(),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
		}
	)
	return a.service.Scan(in)
}

// Random fetches a random quote from the database
func (a *API) Random() *types.QuoteRes {
	var (
		res   = types.NewQuoteRes()
		quote = types.NewQuote()
	)

	scan, err := a.scan()
	logger.InfoOrErr(err).Interface("out", scan).Msg("Requested Quote Batch")
	if *scan.Count == 0 {
		err = errors.New("Error: No quotes")
	}
	if err != nil {
		return res.WithErr(err)
	}

	ceiling := int(math.Min(float64(*scan.Count), float64(a.limit)))
	get, err := a.get(scan.Items[rand.Intn(ceiling)])
	logger.InfoOrErr(err).Interface("out", get).Msg("Requested Quote")
	if err != nil {
		return res.WithErr(err)
	}

	if err = dynamodbattribute.UnmarshalMap(get.Item, quote); err != nil {
		return res.WithErr(err)
	}
	return res.WithQuote(quote)
}
