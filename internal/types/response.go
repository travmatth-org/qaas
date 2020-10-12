package types

// QuoteResponse forms the json response containing a quote
type QuoteResponse struct {
	Quote *Quote
	Err   error
}

// NewQuoteResponse returns a new QuoteResponse
func NewQuoteResponse() *QuoteResponse {
	return &QuoteResponse{}
}

// WithQuote inserts a quote into the QuoteResponse
func (q *QuoteResponse) WithQuote(quote *Quote) *QuoteResponse {
	q.Quote = quote
	return q
}

// WithErr inserts an error into the QuoteResponse
func (q *QuoteResponse) WithErr(err error) *QuoteResponse {
	q.Err = err
	return q
}

// MultiQuoteResponse forms the json response containing multiple quotes
type MultiQuoteResponse struct {
	Quotes []*Quote
	Next   string
	Err    error
}

// NewMultiQuoteResponse returns a new MultiQuoteResponse
func NewMultiQuoteResponse() *MultiQuoteResponse {
	return &MultiQuoteResponse{}
}

// WithQuotes inserts quotes into the MultiQuoteResponse
func (m *MultiQuoteResponse) WithQuotes(q []*Quote) *MultiQuoteResponse {
	m.Quotes = q
	return m
}

// WithNext inserts a next token into the MultiQuoteResponse
func (m *MultiQuoteResponse) WithNext(n string) *MultiQuoteResponse {
	m.Next = n
	return m
}

// WithErr inserts an error into the MultiQuoteResponse
func (m *MultiQuoteResponse) WithErr(e error) *MultiQuoteResponse {
	m.Err = e
	return m
}
