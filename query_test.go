package main

import "testing"

func TestQuery(t *testing.T) {

	testCases := []struct {
		query gmailQuery
		exp   string
	}{
		{
			query: NewerThan(5, day),
			exp:   "newer_than:5d",
		},
		{
			query: NewerThan(5, day).And(From("Alice")),
			exp:   "(newer_than:5d) AND (from:Alice)",
		},
		{
			query: (NewerThan(5, day).And(From("Alice")).Or(From("Bob"))),
			exp:   "((newer_than:5d) AND (from:Alice)) OR (from:Bob)",
		},
		{
			query: NewerThan(5, day).And(From("Alice").Or(From("Bob"))),
			exp:   "(newer_than:5d) AND ((from:Alice) OR (from:Bob))",
		},
	}

	for _, tc := range testCases {
		if tc.query.query != tc.exp {
			t.Fatalf("Expected %s, but got %s", tc.exp, tc.query.query)
		}
	}
}
