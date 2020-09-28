package types

type QuoteResponse struct {
	Quote *Quote
	Err   error
}

type quoteResponseOpt func(q *QuoteResponse) *QuoteResponse

func NewQuoteResponse(opts ...quoteResponseOpt) *QuoteResponse {
	q := &QuoteResponse{}
	for _, opt := range opts {
		q = opt(q)
	}
	return q
}

func (q *QuoteResponse) WithQuote(quote *Quote) *quoteResponseOpt {
	return func(q *QuoteResponse) *QuoteResponse {
		q.Quote = quote
		return q
	}
}

func (q *QuoteResponse) WithErr(err error) *quoteResponseOpt {
	return func(q *QuoteResponse) *QuoteResponse {
		q.Err = err
		return q
	}
}

type MultiQuoteResponse struct {
	Quotes []*Quote
	Next   string
	Err    error
}

type multiQuoteResponseOpt func(q *MultiQuoteResponse) *MultiQuoteResponse

func MultiNewQuoteResponse(opts ...multiQuoteResponseOpt) *MultiQuoteResponse {
	m := &MultiQuoteResponse{}
	for _, opt := range opts {
		m = opt(m)
	}
	return m
}

func (m *MultiQuoteResponse) WithQuotes(q []*Quote) *multiQuoteResponseOpt {
	return func(m *MultiQuoteResponse) *MultiQuoteResponse {
		m.Quotes = q
		return m
	}
}

func (m *MultiQuoteResponse) WithNext(n string) *multiQuoteResponseOpt {
	return func(m *MultiQuoteResponse) *MultiQuoteResponse {
		m.Next = n
		return m
	}
}

func (m *MultiQuoteResponse) WithErr(e error) *multiQuoteResponseOpt {
	return func(m *MultiQuoteResponse) *MultiQuoteResponse {
		m.Err = e
		return m
	}
}