package service

import (
	"fmt"
	"time"
)

type Interval int

const (
	Day Interval = iota
	Month
	Year
)

func (i Interval) String() string {
	return []string{"d", "m", "y"}[i]
}

// Gmail query language
// https://support.google.com/mail/answer/7190?hl=en&co=GENIE.Platform%3DDesktop&oco=0
type GmailQuery struct {
	Query string
}

func (a GmailQuery) And(b GmailQuery) GmailQuery {
	q := fmt.Sprintf("(%s) AND (%s)", a.Query, b.Query)
	return GmailQuery{q}
}

func (a GmailQuery) Or(b GmailQuery) GmailQuery {
	q := fmt.Sprintf("(%s) OR (%s)", a.Query, b.Query)
	return GmailQuery{
		Query: q,
	}
}

func NewerThan(n int, i Interval) GmailQuery {
	query := fmt.Sprintf("newer_than:%d%s", n, i)
	return GmailQuery{query}
}

func OlderThan(n int, i Interval) GmailQuery {
	query := fmt.Sprintf("older_than:%d%s", n, i)
	return GmailQuery{query}
}

func From(id string) GmailQuery {
	query := fmt.Sprintf("from:%s", id)
	return GmailQuery{query}
}

func gmailTime(t time.Time) string {
	return t.Format("04/16/2004")
}
