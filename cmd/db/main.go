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
		config.WithConfigFile(afs.Open),
		config.WithUpdates(os.Args[1:]))
	if err != nil {
		logger.Error().Msg("Error creating config")
		return 1
	}

	// Create API
	a, err := api.New(
		api.WithRegion(c.AWS.Region),
		api.WithSession,
		api.WithEC2(config.IsProd(c)),
		api.WithXray(config.IsProd(c)),
		api.WithNewDDB(c))
	if err != nil {
		logger.Error().Msg("Error configuring API")
		return 1
	}

	var (
		n       = gofakeit.Name()
		t       = []string{gofakeit.Word(), gofakeit.Word()}
		ids     = []string{}
		authors = []string{}
	)

	// Generate & Put new quotes
	var ()

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

		ids = append(ids, quote.ID)
		authors = append(authors, author.Name)
		// records = append(records, author, topics...)
		if err := a.PutNewQuote(quote, author, topics); err != nil {
			fmt.Println(err)
			return 1
		}
	}

	for _, id := range ids {
		in := a.DDB.NewGetQuoteInputByID(id)
		out, err := a.DDB.GetObject(in)
		if err != nil {
			fmt.Println(err)
			return 1
		}
		q, err := a.DDB.QuoteFromObject(out)
		if err != nil {
			fmt.Println(err)
			return 1
		}
		fmt.Printf("%+v\n", q)
	}

	for _, topic := range t {
		var (
			last *types.Record = nil
		)
		for {
			res := a.GetQuotesByTopic(topic, last)
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
			res := a.GetQuotesByAuthor(author, last)
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
	quote := a.GetRandomQuote()
	fmt.Printf("%+v", quote)
	return 0
}

func main() {
	os.Exit(start())
}
