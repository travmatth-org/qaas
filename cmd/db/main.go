package main

import (
	"fmt"
	"os"

	"github.com/brianvoe/gofakeit/v5"

	"github.com/travmatth-org/qaas/internal/afs"
	"github.com/travmatth-org/qaas/internal/api"
	"github.com/travmatth-org/qaas/internal/config"
	"github.com/travmatth-org/qaas/internal/logger"
	"github.com/travmatth-org/qaas/internal/types"
)

func start() int {
	// init filesystem
	afs := afs.New().WithCachedFs()
	// Load config options
	c, err := config.New(
		config.WithFile(afs.Open),
		config.Update(os.Args[1:]))
	if err != nil {
		logger.Error().Msg("Error creating config")
		return 1
	}

	// Create API
	a, err := api.New(
		api.WithSession,
		api.WithEC2(config.IsProd(c)),
		api.WithXray(config.IsProd(c)),
		api.WithDynamoDBService(c))
	if err != nil {
		logger.Error().Msg("Error configuring API")
		return 1
	}

	var (
		n       = gofakeit.Name()
		t       = []string{gofakeit.Word(), gofakeit.Word()}
		authors = []string{}
	)

	for i := 0; i < 10; i++ {
		var (
			q      = gofakeit.Phrase()
			quote  = types.NewQuote().NewID().WithText(q).WithAuthor(n).WithTopics(t)
			author = types.NewRecord().WithName(n).WithQuoteID(quote.ID)
			topics = make([]*types.Record, 0)
		)

		for _, _t := range t {
			topic := types.NewRecord().WithName(_t).WithQuoteID(quote.ID)
			topics = append(topics, topic)
		}

		authors = append(authors, author.Name)
		// records = append(records, author, topics...)
		err := a.Put(
			a.PutWithQuote(quote),
			a.PutWithAuthor(author),
			a.PutWithTopics(topics))
		if err != nil {
			fmt.Println(err)
			return 1
		}
	}

	for _, topic := range t {
		var (
			last *types.Record = nil
		)
		for {
			res := a.Get(a.Table.Topic, topic, last)
			if res.Err != nil {
				fmt.Println(res.Err)
				return 1
			}
			for _, q := range res.Quotes {
				fmt.Printf("%+v\n", q)
			}
			if res.Next == nil {
				break
			}
			last = res.Next
		}
	}
	for _, author := range authors {
		var (
			last *types.Record = nil
		)
		for {
			res := a.Get(a.Table.Author, author, last)
			if res.Err != nil {
				fmt.Println(res.Err)
				return 1
			}
			if res.Next == nil {
				break
			}
			last = res.Next
		}
	}
	quote := a.Random()
	fmt.Printf("%+v", quote)
	return 0
}

func main() {
	os.Exit(start())
}
