package dumbql_test

import (
	"fmt"

	"github.com/defer-panic/dumbql/query"
)

type User struct {
	ID       int64   `dumbql:"id"`
	Name     string  `dumbql:"name"`
	Age      int64   `dumbql:"age"`
	Score    float64 `dumbql:"score"`
	Location string  `dumbql:"location"`
	Role     string  `dumbql:"role"`
}

func Example_simpleMatching() {
	user := &User{
		ID:       1,
		Name:     "John Doe",
		Age:      30,
		Score:    4.5,
		Location: "New York",
		Role:     "admin",
	}

	// Parse a simple equality query
	q := `name = "John Doe"`
	ast, _ := query.Parse("test", []byte(q))
	expr := ast.(query.Expr)

	// Create a matcher
	matcher := &query.DefaultMatcher{}
	result := expr.Match(user, matcher)

	fmt.Printf("%s: %v\n", q, result)
	// Output: name = "John Doe": true
}

func Example_complexMatching() {
	user := &User{
		ID:       1,
		Name:     "John Doe",
		Age:      30,
		Score:    4.5,
		Location: "New York",
		Role:     "admin",
	}

	// Parse a complex query with multiple conditions
	q := `age >= 25 and location:["New York", "Los Angeles"] and score > 4.0`
	ast, _ := query.Parse("test", []byte(q))
	expr := ast.(query.Expr)

	// Create a matcher
	matcher := &query.DefaultMatcher{}
	result := expr.Match(user, matcher)

	fmt.Printf("%s: %v\n", q, result)
	// Output: age >= 25 and location:["New York", "Los Angeles"] and score > 4.0: true
}

func Example_numericComparisons() {
	user := &User{
		ID:       1,
		Name:     "John Doe",
		Age:      30,
		Score:    4.5,
		Location: "New York",
		Role:     "admin",
	}

	// Test various numeric comparisons
	queries := []string{
		`age > 20`,
		`age < 40`,
		`age >= 30`,
		`age <= 30`,
		`score > 4.0`,
		`score < 5.0`,
	}

	matcher := &query.DefaultMatcher{}

	for _, q := range queries {
		ast, _ := query.Parse("test", []byte(q))
		expr := ast.(query.Expr)
		result := expr.Match(user, matcher)
		fmt.Printf("Query '%s' match result: %v\n", q, result)
	}
	// Output: Query 'age > 20' match result: true
	// Query 'age < 40' match result: true
	// Query 'age >= 30' match result: true
	// Query 'age <= 30' match result: true
	// Query 'score > 4.0' match result: true
	// Query 'score < 5.0' match result: true
}

func Example_stringOperations() {
	user := &User{
		ID:       1,
		Name:     "John Doe",
		Age:      30,
		Score:    4.5,
		Location: "New York",
		Role:     "admin",
	}

	// Test various string operations
	queries := []string{
		`name:"John Doe"`,
		`name~"John"`,
		`location:"New York"`,
		`role:admin`,
	}

	matcher := &query.DefaultMatcher{}

	for _, q := range queries {
		ast, _ := query.Parse("test", []byte(q))
		expr := ast.(query.Expr)
		result := expr.Match(user, matcher)
		fmt.Printf("Query '%s' match result: %v\n", q, result)
	}
	// Output:
	// Query 'name:"John Doe"' match result: true
	// Query 'name~"John"' match result: true
	// Query 'location:"New York"' match result: true
	// Query 'role:admin' match result: true
}

func Example_notExpressions() {
	user := &User{
		ID:       1,
		Name:     "John Doe",
		Age:      30,
		Score:    4.5,
		Location: "New York",
		Role:     "admin",
	}

	// Test NOT expressions
	queries := []string{
		`not age < 25`,
		`not location:"Los Angeles"`,
		`not (role:"user" and score < 3.0)`,
	}

	matcher := &query.DefaultMatcher{}

	for _, q := range queries {
		ast, _ := query.Parse("test", []byte(q))
		expr := ast.(query.Expr)
		result := expr.Match(user, matcher)
		fmt.Printf("Query '%s' match result: %v\n", q, result)
	}
	// Output:
	// Query 'not age < 25' match result: true
	// Query 'not location:"Los Angeles"' match result: true
	// Query 'not (role:"user" and score < 3.0)' match result: true
}

func Example_multiMatch() {
	users := []User{
		{
			ID:       1,
			Name:     "John Doe",
			Age:      30,
			Score:    4.5,
			Location: "New York",
			Role:     "admin",
		},
		{
			ID:       2,
			Name:     "Jane Smith",
			Age:      25,
			Score:    3.8,
			Location: "Los Angeles",
			Role:     "user",
		},
		{
			ID:       3,
			Name:     "Bob Johnson",
			Age:      35,
			Score:    4.2,
			Location: "Chicago",
			Role:     "user",
		},
		// This one will be dropped:
		{
			ID:       4,
			Name:     "Alice Smith",
			Age:      25,
			Score:    3.8,
			Location: "Los Angeles",
			Role:     "admin",
		},
	}

	q := `(age >= 30 and score > 4.0) or (location:"Los Angeles" and role:"user")`
	ast, _ := query.Parse("test", []byte(q))
	expr := ast.(query.Expr)

	matcher := &query.DefaultMatcher{}

	filtered := make([]User, 0, len(users))

	for _, user := range users {
		if expr.Match(&user, matcher) {
			filtered = append(filtered, user)
		}
	}

	fmt.Println(filtered)
	// Output:
	// [{1 John Doe 30 4.5 New York admin} {2 Jane Smith 25 3.8 Los Angeles user} {3 Bob Johnson 35 4.2 Chicago user}]
}

func Example_oneOfExpression() {
	user := &User{
		ID:       1,
		Name:     "John Doe",
		Age:      30,
		Score:    4.5,
		Location: "New York",
		Role:     "admin",
	}

	// Test OneOf expressions
	queries := []string{
		`location:["New York", "Los Angeles", "Chicago"]`,
		`role:["admin", "superuser"]`,
		`age:[25, 30, 35]`,
	}

	matcher := &query.DefaultMatcher{}

	for _, q := range queries {
		ast, _ := query.Parse("test", []byte(q))
		expr := ast.(query.Expr)
		result := expr.Match(user, matcher)
		fmt.Printf("Query '%s' match result: %v\n", q, result)
	}
	// Output:
	// Query 'location:["New York", "Los Angeles", "Chicago"]' match result: true
	// Query 'role:["admin", "superuser"]' match result: true
	// Query 'age:[25, 30, 35]' match result: true
}

func Example_edgeCases() {
	user := &User{
		ID:       1,
		Name:     "John Doe",
		Age:      30,
		Score:    4.5,
		Location: "New York",
		Role:     "admin",
	}
	// Test edge cases and special scenarios
	queries := []string{
		// Non-existent field
		`nonexistent:"value"`,
		// Invalid type comparison
		`age:"not a number"`,
		// Empty string matching
		`name:""`,
		// Zero value matching
		`score:0`,
		// Complex nested expression
		`(age > 20 and age < 40) and (score >= 4.0 or role:"admin")`,
	}

	matcher := &query.DefaultMatcher{}

	for _, q := range queries {
		ast, _ := query.Parse("test", []byte(q))
		expr := ast.(query.Expr)
		result := expr.Match(user, matcher)
		fmt.Printf("Query '%s' match result: %v\n", q, result)
	}

	// Output:
	// Query 'nonexistent:"value"' match result: false
	// Query 'age:"not a number"' match result: false
	// Query 'name:""' match result: false
	// Query 'score:0' match result: false
	// Query '(age > 20 and age < 40) and (score >= 4.0 or role:"admin")' match result: true
}
