package main

import "fmt"

type Interval int

const (
	day Interval = iota
	month
	year
)

func (i Interval) String() string {
	return []string{"d", "m", "y"}[i]
}

// Gmail query language
// https://support.google.com/mail/answer/7190?hl=en&co=GENIE.Platform%3DDesktop&oco=0
type gmailQuery struct {
	query string
}

func (a gmailQuery) And(b gmailQuery) gmailQuery {
	q := fmt.Sprintf("(%s) AND (%s)", a.query, b.query)
	return gmailQuery{q}
}

func (a gmailQuery) Or(b gmailQuery) gmailQuery {
	q := fmt.Sprintf("(%s) OR (%s)", a.query, b.query)
	return gmailQuery{
		query: q,
	}
}

func NewerThan(n int, i Interval) gmailQuery {
	query := fmt.Sprintf("newer_than:%d%s", n, i)
	return gmailQuery{query}
}

func OlderThan(n int, i Interval) gmailQuery {
	query := fmt.Sprintf("older_than:%d%s", n, i)
	return gmailQuery{query}
}

func From(id string) gmailQuery {
	query := fmt.Sprintf("from:%s", id)
	return gmailQuery{query}
}
