package service

import "testing"

func TestQuery(t *testing.T) {

	testCases := []struct {
		query GmailQuery
		exp   string
	}{
		{
			query: NewerThan(5, Day),
			exp:   "newer_than:5d",
		},
		{
			query: NewerThan(5, Day).And(From("Alice")),
			exp:   "(newer_than:5d) AND (from:Alice)",
		},
		{
			query: (NewerThan(5, Day).And(From("Alice")).Or(From("Bob"))),
			exp:   "((newer_than:5d) AND (from:Alice)) OR (from:Bob)",
		},
		{
			query: NewerThan(5, Day).And(From("Alice").Or(From("Bob"))),
			exp:   "(newer_than:5d) AND ((from:Alice) OR (from:Bob))",
		},
	}

	for _, tc := range testCases {
		if tc.query.Query != tc.exp {
			t.Fatalf("Expected %s, but got %s", tc.exp, tc.query.Query)
		}
	}
}
