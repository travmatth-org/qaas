package types

type QuoteResponse struct {
	Quote *Quote
	Err   error
}

func NewQuoteResponse() *QuoteResponse {
	return &QuoteResponse{}
}

func (q *QuoteResponse) WithQuote(quote *Quote) *QuoteResponse {
	q.Quote = quote
	return q
}

func (q *QuoteResponse) WithErr(err error) *QuoteResponse {
	q.Err = err
	return q
}

type MultiQuoteResponse struct {
	Quotes []*Quote
	Next   string
	Err    error
}

func NewMultiQuoteResponse() *MultiQuoteResponse {
	return &MultiQuoteResponse{}
}

func (m *MultiQuoteResponse) WithQuotes(q []*Quote) *MultiQuoteResponse {
	m.Quotes = q
	return m
}

func (m *MultiQuoteResponse) WithNext(n string) *MultiQuoteResponse {
	m.Next = n
	return m
}

func (m *MultiQuoteResponse) WithErr(e error) *MultiQuoteResponse {
	m.Err = e
	return m
}