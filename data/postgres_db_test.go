package data

import (
	"testing"
)

func TestGetSQLForListAll(t *testing.T) {
	filters := make(map[string][]string)
	res, args := getSQLSelect(filters)
	if res != "SELECT * FROM users" && len(args) != 0 {
		t.Fatalf("No filters --- SQL generation for get list all is failed")
	}

	filters["first_name"] = []string{"Roman", "Jan"}
	filters["country"] = []string{"Russia"}
	res, args = getSQLSelect(filters)
	expected := "SELECT * FROM users WHERE first_name IN (%1, %2) AND country = %3;"
	if len(res) != len(expected) && len(args) != 3 {
		t.Fatalf("SQL generation for get list all is failed\n Expected: %s . \n Got: %s", expected, res)
	}
}
