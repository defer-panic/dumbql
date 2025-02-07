package dumbql_test

import (
	"testing"

	"github.com/defer-panic/dumbql"
	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		// Simple field expression.
		{
			input: "status:200",
			want:  "(= status 200)",
		},
		// Floating-point number.
		{
			input: "eps<0.003",
			want:  "(< eps 0.003000)",
		},
		// Using <= operator.
		{
			input: "eps<=0.003",
			want:  "(<= eps 0.003000)",
		},
		// Using >= operator.
		{
			input: "eps>=0.003",
			want:  "(>= eps 0.003000)",
		},
		// Using > operator.
		{
			input: "eps>0.003",
			want:  "(> eps 0.003000)",
		},
		// Using not-equals with !: operator.
		{
			input: "eps!:0.003",
			want:  "(!= eps 0.003000)",
		},
		// Combined with AND.
		{
			input: "status:200 and eps < 0.003",
			want:  "(and (= status 200) (< eps 0.003000))",
		},
		// Combined with OR.
		{
			input: "status:200 or eps<0.003",
			want:  "(or (= status 200) (< eps 0.003000))",
		},
		// Mixed operators: AND with not-equals.
		{
			input: "status:200 and eps!=0.003",
			want:  "(and (= status 200) (!= eps 0.003000))",
		},
		// Nested parentheses.
		{
			input: "((status:200))",
			want:  "(= status 200)",
		},
		// Extra whitespace.
		{
			input: "   status  :   200    and   eps  <  0.003   ",
			want:  "(and (= status 200) (< eps 0.003000))",
		},
		// Uppercase boolean operator.
		{
			input: "status:200 AND eps<0.003",
			want:  "(and (= status 200) (< eps 0.003000))",
		},
		// Array literal in a field expression.
		{
			input: "req.fields.ext:[\"jpg\", \"png\"]",
			want:  "(= req.fields.ext [\"jpg\" \"png\"])",
		},
		// Array with a single element.
		{
			input: "tags:[\"urgent\"]",
			want:  "(= tags [\"urgent\"])",
		},
		// Empty array literal.
		{
			input: "tags:[]",
			want:  "(= tags [])",
		},
		// A complex expression combining several constructs.
		{
			input: "status : 200 and eps < 0.003 and (req.fields.ext:[\"jpg\", \"png\"])",
			want:  "(and (and (= status 200) (< eps 0.003000)) (= req.fields.ext [\"jpg\" \"png\"]))",
		},
		// NOT with parentheses.
		{
			input: "not (status:200)",
			want:  "(not (= status 200))",
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			ast, err := dumbql.Parse("input", []byte(test.input))
			require.NoError(t, err, "parsing error for input: %s", test.input)

			require.Equal(t, test.want, ast.(dumbql.Expr).String())
		})
	}
}
