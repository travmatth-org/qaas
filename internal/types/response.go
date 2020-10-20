package types

import (
	"encoding/json"
)

// QuoteRes forms the json response containing a quote
type QuoteRes struct {
	Type  string
	Quote *Quote
	Err   error `json:"error"`
}

// NewQuoteRes returns a new QuoteRes
func NewQuoteRes() *QuoteRes {
	return &QuoteRes{
		Type:  "Quote",
		Quote: nil,
		Err:   nil,
	}
}

// WithQuote inserts a quote into the QuoteRes
func (q *QuoteRes) WithQuote(quote *Quote) *QuoteRes {
	q.Quote = quote
	return q
}

// WithErr inserts an error into the QuoteRes
func (q *QuoteRes) WithErr(err error) *QuoteRes {
	q.Err = err
	return q
}

// JSON return a byte slice containing the JSON representation of the QuoteRes
func (q *QuoteRes) JSON() []byte {
	bytes, _ := json.Marshal(q)
	return bytes
}

// Error returns the string representation of the QuoteRes error field
func (q *QuoteRes) Error() string {
	if q.Err != nil {
		return q.Err.Error()
	}
	return ""
}

// MarshalZerologObject logs the given object to zerolog
func (q *QuoteRes) MarshalZerologObject(e *ZLEvent) {
	if q.Quote != nil {
		e.Str("ID", q.Quote.ID)
	} else {
		e.Str("ID", "Empty Quote")
	}
	if q.Err != nil {
		e.Str("err", q.Err.Error())
	}
}

// MultiQuoteRes forms the json response containing multiple quotes
type MultiQuoteRes struct {
	Type   string
	Quotes []*Quote
	Next   *Record
	Err    error `json:"error"`
}

// NewMultiQuoteRes returns a new MultiQuoteRes
func NewMultiQuoteRes() *MultiQuoteRes {
	return &MultiQuoteRes{
		Type:   "multi",
		Quotes: []*Quote{},
		Next:   nil,
		Err:    nil,
	}
}

// WithQuotes inserts quotes into the MultiQuoteRes
func (m *MultiQuoteRes) WithQuotes(q []*Quote) *MultiQuoteRes {
	m.Quotes = q
	return m
}

// WithNext inserts a next token into the MultiQuoteRes
func (m *MultiQuoteRes) WithNext(r *Record) *MultiQuoteRes {
	m.Next = r
	return m
}

// WithErr inserts an error into the MultiQuoteRes
func (m *MultiQuoteRes) WithErr(e error) *MultiQuoteRes {
	m.Err = e
	return m
}

// JSON return a byte slice containing the JSON representation of the MultiQuoteRes
func (m *MultiQuoteRes) JSON() []byte {
	bytes, _ := json.Marshal(m)
	return bytes
}

// MarshalZerologObject logs the given object to zerolog
func (m *MultiQuoteRes) MarshalZerologObject(e *ZLEvent) {
	ids := []string{}
	if m.Quotes != nil {
		for _, q := range m.Quotes {
			ids = append(ids, q.ID)
		}
	}
	e.Strs("ids", ids)
	if m.Err != nil {
		e.Str("err", m.Err.Error())
	}
	e.Interface("next", m.Next)
}

// Error returns the string representation of the QuoteRes error field
func (m *MultiQuoteRes) Error() string {
	if m.Err != nil {
		return m.Err.Error()
	}
	return ""
}

// GetBatch represents the client request for a batch of quotes with attribute
type GetBatch struct {
	Name  string `validate:"regexp=^(author|topic)$"`
	Value string `validate:"min=3,max=20,regexp=^[a-zA-Z0-9 ]*$"`
	Start *Record
}

func NewGetBatch() *GetBatch {
	return &GetBatch{}
}
