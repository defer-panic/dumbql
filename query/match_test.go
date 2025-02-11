package query_test

import (
	"testing"

	"github.com/defer-panic/dumbql/query"
	"github.com/stretchr/testify/assert"
)

type person struct {
	Name     string  `dumbql:"name"`
	Age      int64   `dumbql:"age"`
	Height   float64 `dumbql:"height"`
	IsMember bool
}

func TestDefaultMatcher_MatchAnd(t *testing.T) { //nolint:funlen
	matcher := &query.DefaultMatcher{}
	target := person{Name: "John", Age: 30}

	tests := []struct {
		name  string
		left  query.Expr
		right query.Expr
		want  bool
	}{
		{
			name: "both conditions true",
			left: &query.FieldExpr{
				Field: "name",
				Op:    query.Equal,
				Value: &query.StringLiteral{StringValue: "John"},
			},
			right: &query.FieldExpr{
				Field: "age",
				Op:    query.Equal,
				Value: &query.IntegerLiteral{IntegerValue: 30},
			},
			want: true,
		},
		{
			name: "left condition false",
			left: &query.FieldExpr{
				Field: "name",
				Op:    query.Equal,
				Value: &query.StringLiteral{StringValue: "Jane"},
			},
			right: &query.FieldExpr{
				Field: "age",
				Op:    query.Equal,
				Value: &query.IntegerLiteral{IntegerValue: 30},
			},
			want: false,
		},
		{
			name: "right condition false",
			left: &query.FieldExpr{
				Field: "name",
				Op:    query.Equal,
				Value: &query.StringLiteral{StringValue: "John"},
			},
			right: &query.FieldExpr{
				Field: "age",
				Op:    query.Equal,
				Value: &query.IntegerLiteral{IntegerValue: 25},
			},
			want: false,
		},
		{
			name: "both conditions false",
			left: &query.FieldExpr{
				Field: "name",
				Op:    query.Equal,
				Value: &query.StringLiteral{StringValue: "Jane"},
			},
			right: &query.FieldExpr{
				Field: "age",
				Op:    query.Equal,
				Value: &query.IntegerLiteral{IntegerValue: 25},
			},
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := matcher.MatchAnd(target, test.left, test.right)
			assert.Equal(t, test.want, result)
		})
	}
}

func TestDefaultMatcher_MatchOr(t *testing.T) { //nolint:funlen
	matcher := &query.DefaultMatcher{}
	target := person{Name: "John", Age: 30}

	tests := []struct {
		name  string
		left  query.Expr
		right query.Expr
		want  bool
	}{
		{
			name: "both conditions true",
			left: &query.FieldExpr{
				Field: "name",
				Op:    query.Equal,
				Value: &query.StringLiteral{StringValue: "John"},
			},
			right: &query.FieldExpr{
				Field: "age",
				Op:    query.Equal,
				Value: &query.IntegerLiteral{IntegerValue: 30},
			},
			want: true,
		},
		{
			name: "left condition true only",
			left: &query.FieldExpr{
				Field: "name",
				Op:    query.Equal,
				Value: &query.StringLiteral{StringValue: "John"},
			},
			right: &query.FieldExpr{
				Field: "age",
				Op:    query.Equal,
				Value: &query.IntegerLiteral{IntegerValue: 25},
			},
			want: true,
		},
		{
			name: "right condition true only",
			left: &query.FieldExpr{
				Field: "name",
				Op:    query.Equal,
				Value: &query.StringLiteral{StringValue: "Jane"},
			},
			right: &query.FieldExpr{
				Field: "age",
				Op:    query.Equal,
				Value: &query.IntegerLiteral{IntegerValue: 30},
			},
			want: true,
		},
		{
			name: "both conditions false",
			left: &query.FieldExpr{
				Field: "name",
				Op:    query.Equal,
				Value: &query.StringLiteral{StringValue: "Jane"},
			},
			right: &query.FieldExpr{
				Field: "age",
				Op:    query.Equal,
				Value: &query.IntegerLiteral{IntegerValue: 25},
			},
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := matcher.MatchOr(target, test.left, test.right)
			assert.Equal(t, test.want, result)
		})
	}
}

func TestDefaultMatcher_MatchNot(t *testing.T) {
	matcher := &query.DefaultMatcher{}
	target := person{Name: "John", Age: 30}

	tests := []struct {
		name string
		expr query.Expr
		want bool
	}{
		{
			name: "negate true condition",
			expr: &query.FieldExpr{
				Field: "name",
				Op:    query.Equal,
				Value: &query.StringLiteral{StringValue: "John"},
			},
			want: false,
		},
		{
			name: "negate false condition",
			expr: &query.FieldExpr{
				Field: "name",
				Op:    query.Equal,
				Value: &query.StringLiteral{StringValue: "Jane"},
			},
			want: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := matcher.MatchNot(target, test.expr)
			assert.Equal(t, test.want, result)
		})
	}
}

func TestDefaultMatcher_MatchField(t *testing.T) {
	matcher := &query.DefaultMatcher{}
	target := person{
		Name:     "John",
		Age:      30,
		Height:   1.75,
		IsMember: true,
	}

	tests := []struct {
		name  string
		field string
		value query.Valuer
		op    query.FieldOperator
		want  bool
	}{
		{
			name:  "string equal match",
			field: "name",
			value: &query.StringLiteral{StringValue: "John"},
			op:    query.Equal,
			want:  true,
		},
		{
			name:  "string not equal match",
			field: "name",
			value: &query.StringLiteral{StringValue: "Jane"},
			op:    query.NotEqual,
			want:  true,
		},
		{
			name:  "integer equal match",
			field: "age",
			value: &query.IntegerLiteral{IntegerValue: 30},
			op:    query.Equal,
			want:  true,
		},
		{
			name:  "float greater than match",
			field: "height",
			value: &query.NumberLiteral{NumberValue: 1.70},
			op:    query.GreaterThan,
			want:  true,
		},
		{
			name:  "non-existent field",
			field: "invalid",
			value: &query.StringLiteral{StringValue: "test"},
			op:    query.Equal,
			want:  false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := matcher.MatchField(target, test.field, test.value, test.op)
			assert.Equal(t, test.want, result)
		})
	}
}

func TestDefaultMatcher_MatchValue(t *testing.T) {
	t.Run("string", testMatchValueString)
	t.Run("integer", testMatchValueInteger)
	t.Run("float", testMatchValueFloat)
	t.Run("type mismatch", testMatchValueTypeMismatch)
}

func testMatchValueString(t *testing.T) { //nolint:funlen
	matcher := &query.DefaultMatcher{}
	tests := []struct {
		name   string
		target any
		value  query.Valuer
		op     query.FieldOperator
		want   bool
	}{
		{
			name:   "equal - match",
			target: "hello",
			value:  &query.StringLiteral{StringValue: "hello"},
			op:     query.Equal,
			want:   true,
		},
		{
			name:   "equal - no match",
			target: "hello",
			value:  &query.StringLiteral{StringValue: "world"},
			op:     query.Equal,
			want:   false,
		},
		{
			name:   "not equal - match",
			target: "hello",
			value:  &query.StringLiteral{StringValue: "world"},
			op:     query.NotEqual,
			want:   true,
		},
		{
			name:   "not equal - no match",
			target: "hello",
			value:  &query.StringLiteral{StringValue: "hello"},
			op:     query.NotEqual,
			want:   false,
		},
		{
			name:   "like - match",
			target: "hello world",
			value:  &query.StringLiteral{StringValue: "world"},
			op:     query.Like,
			want:   true,
		},
		{
			name:   "like - no match",
			target: "hello world",
			value:  &query.StringLiteral{StringValue: "universe"},
			op:     query.Like,
			want:   false,
		},
		{
			name:   "greater than - invalid",
			target: "hello",
			value:  &query.StringLiteral{StringValue: "world"},
			op:     query.GreaterThan,
			want:   false,
		},
		{
			name:   "greater than or equal - invalid",
			target: "hello",
			value:  &query.StringLiteral{StringValue: "world"},
			op:     query.GreaterThanOrEqual,
			want:   false,
		},
		{
			name:   "less than - invalid",
			target: "hello",
			value:  &query.StringLiteral{StringValue: "world"},
			op:     query.LessThan,
			want:   false,
		},
		{
			name:   "less than or equal - invalid",
			target: "hello",
			value:  &query.StringLiteral{StringValue: "world"},
			op:     query.LessThanOrEqual,
			want:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := matcher.MatchValue(test.target, test.value, test.op)
			assert.Equal(t, test.want, result)
		})
	}
}

func testMatchValueInteger(t *testing.T) { //nolint:funlen
	matcher := &query.DefaultMatcher{}
	tests := []struct {
		name   string
		target any
		value  query.Valuer
		op     query.FieldOperator
		want   bool
	}{
		{
			name:   "equal - match",
			target: int64(42),
			value:  &query.IntegerLiteral{IntegerValue: 42},
			op:     query.Equal,
			want:   true,
		},
		{
			name:   "equal - no match",
			target: int64(42),
			value:  &query.IntegerLiteral{IntegerValue: 24},
			op:     query.Equal,
			want:   false,
		},
		{
			name:   "not equal - match",
			target: int64(42),
			value:  &query.IntegerLiteral{IntegerValue: 24},
			op:     query.NotEqual,
			want:   true,
		},
		{
			name:   "not equal - no match",
			target: int64(42),
			value:  &query.IntegerLiteral{IntegerValue: 42},
			op:     query.NotEqual,
			want:   false,
		},
		{
			name:   "greater than - match",
			target: int64(42),
			value:  &query.IntegerLiteral{IntegerValue: 24},
			op:     query.GreaterThan,
			want:   true,
		},
		{
			name:   "greater than - no match",
			target: int64(24),
			value:  &query.IntegerLiteral{IntegerValue: 42},
			op:     query.GreaterThan,
			want:   false,
		},
		{
			name:   "greater than or equal - match (greater)",
			target: int64(42),
			value:  &query.IntegerLiteral{IntegerValue: 24},
			op:     query.GreaterThanOrEqual,
			want:   true,
		},
		{
			name:   "greater than or equal - match (equal)",
			target: int64(42),
			value:  &query.IntegerLiteral{IntegerValue: 42},
			op:     query.GreaterThanOrEqual,
			want:   true,
		},
		{
			name:   "greater than or equal - no match",
			target: int64(24),
			value:  &query.IntegerLiteral{IntegerValue: 42},
			op:     query.GreaterThanOrEqual,
			want:   false,
		},
		{
			name:   "less than - match",
			target: int64(24),
			value:  &query.IntegerLiteral{IntegerValue: 42},
			op:     query.LessThan,
			want:   true,
		},
		{
			name:   "less than - no match",
			target: int64(42),
			value:  &query.IntegerLiteral{IntegerValue: 24},
			op:     query.LessThan,
			want:   false,
		},
		{
			name:   "less than or equal - match (less)",
			target: int64(24),
			value:  &query.IntegerLiteral{IntegerValue: 42},
			op:     query.LessThanOrEqual,
			want:   true,
		},
		{
			name:   "less than or equal - match (equal)",
			target: int64(42),
			value:  &query.IntegerLiteral{IntegerValue: 42},
			op:     query.LessThanOrEqual,
			want:   true,
		},
		{
			name:   "less than or equal - no match",
			target: int64(42),
			value:  &query.IntegerLiteral{IntegerValue: 24},
			op:     query.LessThanOrEqual,
			want:   false,
		},
		{
			name:   "like - invalid",
			target: int64(42),
			value:  &query.IntegerLiteral{IntegerValue: 24},
			op:     query.Like,
			want:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := matcher.MatchValue(test.target, test.value, test.op)
			assert.Equal(t, test.want, result)
		})
	}
}

func testMatchValueFloat(t *testing.T) { //nolint:funlen
	matcher := &query.DefaultMatcher{}
	tests := []struct {
		name   string
		target any
		value  query.Valuer
		op     query.FieldOperator
		want   bool
	}{
		{
			name:   "equal - match",
			target: 3.14,
			value:  &query.NumberLiteral{NumberValue: 3.14},
			op:     query.Equal,
			want:   true,
		},
		{
			name:   "equal - no match",
			target: 3.14,
			value:  &query.NumberLiteral{NumberValue: 2.718},
			op:     query.Equal,
			want:   false,
		},
		{
			name:   "not equal - match",
			target: 3.14,
			value:  &query.NumberLiteral{NumberValue: 2.718},
			op:     query.NotEqual,
			want:   true,
		},
		{
			name:   "not equal - no match",
			target: 3.14,
			value:  &query.NumberLiteral{NumberValue: 3.14},
			op:     query.NotEqual,
			want:   false,
		},
		{
			name:   "greater than - match",
			target: 3.14,
			value:  &query.NumberLiteral{NumberValue: 2.718},
			op:     query.GreaterThan,
			want:   true,
		},
		{
			name:   "greater than - no match",
			target: 2.718,
			value:  &query.NumberLiteral{NumberValue: 3.14},
			op:     query.GreaterThan,
			want:   false,
		},
		{
			name:   "greater than or equal - match (greater)",
			target: 3.14,
			value:  &query.NumberLiteral{NumberValue: 2.718},
			op:     query.GreaterThanOrEqual,
			want:   true,
		},
		{
			name:   "greater than or equal - match (equal)",
			target: 3.14,
			value:  &query.NumberLiteral{NumberValue: 3.14},
			op:     query.GreaterThanOrEqual,
			want:   true,
		},
		{
			name:   "greater than or equal - no match",
			target: 2.718,
			value:  &query.NumberLiteral{NumberValue: 3.14},
			op:     query.GreaterThanOrEqual,
			want:   false,
		},
		{
			name:   "less than - match",
			target: 2.718,
			value:  &query.NumberLiteral{NumberValue: 3.14},
			op:     query.LessThan,
			want:   true,
		},
		{
			name:   "less than - no match",
			target: 3.14,
			value:  &query.NumberLiteral{NumberValue: 2.718},
			op:     query.LessThan,
			want:   false,
		},
		{
			name:   "less than or equal - match (less)",
			target: 2.718,
			value:  &query.NumberLiteral{NumberValue: 3.14},
			op:     query.LessThanOrEqual,
			want:   true,
		},
		{
			name:   "less than or equal - match (equal)",
			target: 3.14,
			value:  &query.NumberLiteral{NumberValue: 3.14},
			op:     query.LessThanOrEqual,
			want:   true,
		},
		{
			name:   "less than or equal - no match",
			target: 3.14,
			value:  &query.NumberLiteral{NumberValue: 2.718},
			op:     query.LessThanOrEqual,
			want:   false,
		},
		{
			name:   "like - invalid",
			target: 3.14,
			value:  &query.NumberLiteral{NumberValue: 2.718},
			op:     query.Like,
			want:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := matcher.MatchValue(test.target, test.value, test.op)
			assert.Equal(t, test.want, result)
		})
	}
}

func testMatchValueTypeMismatch(t *testing.T) {
	matcher := &query.DefaultMatcher{}
	tests := []struct {
		name   string
		target any
		value  query.Valuer
		op     query.FieldOperator
		want   bool
	}{
		{
			name:   "string target with integer value",
			target: "42",
			value:  &query.IntegerLiteral{IntegerValue: 42},
			op:     query.Equal,
			want:   false,
		},
		{
			name:   "integer target with float value",
			target: int64(42),
			value:  &query.NumberLiteral{NumberValue: 42.0},
			op:     query.Equal,
			want:   false,
		},
		{
			name:   "float target with string value",
			target: 3.14,
			value:  &query.StringLiteral{StringValue: "3.14"},
			op:     query.Equal,
			want:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := matcher.MatchValue(test.target, test.value, test.op)
			assert.Equal(t, test.want, result)
		})
	}
}

func TestBinaryExpr_Match(t *testing.T) { //nolint:funlen
	target := person{Name: "John", Age: 30}
	matcher := &query.DefaultMatcher{}

	tests := []struct {
		name string
		expr *query.BinaryExpr
		want bool
	}{
		{
			name: "AND - both true",
			expr: &query.BinaryExpr{
				Left: &query.FieldExpr{
					Field: "name",
					Op:    query.Equal,
					Value: &query.StringLiteral{StringValue: "John"},
				},
				Op: query.And,
				Right: &query.FieldExpr{
					Field: "age",
					Op:    query.Equal,
					Value: &query.IntegerLiteral{IntegerValue: 30},
				},
			},
			want: true,
		},
		{
			name: "AND - left false",
			expr: &query.BinaryExpr{
				Left: &query.FieldExpr{
					Field: "name",
					Op:    query.Equal,
					Value: &query.StringLiteral{StringValue: "Jane"},
				},
				Op: query.And,
				Right: &query.FieldExpr{
					Field: "age",
					Op:    query.Equal,
					Value: &query.IntegerLiteral{IntegerValue: 30},
				},
			},
			want: false,
		},
		{
			name: "OR - both true",
			expr: &query.BinaryExpr{
				Left: &query.FieldExpr{
					Field: "name",
					Op:    query.Equal,
					Value: &query.StringLiteral{StringValue: "John"},
				},
				Op: query.Or,
				Right: &query.FieldExpr{
					Field: "age",
					Op:    query.Equal,
					Value: &query.IntegerLiteral{IntegerValue: 30},
				},
			},
			want: true,
		},
		{
			name: "OR - one true",
			expr: &query.BinaryExpr{
				Left: &query.FieldExpr{
					Field: "name",
					Op:    query.Equal,
					Value: &query.StringLiteral{StringValue: "John"},
				},
				Op: query.Or,
				Right: &query.FieldExpr{
					Field: "age",
					Op:    query.Equal,
					Value: &query.IntegerLiteral{IntegerValue: 25},
				},
			},
			want: true,
		},
		{
			name: "OR - both false",
			expr: &query.BinaryExpr{
				Left: &query.FieldExpr{
					Field: "name",
					Op:    query.Equal,
					Value: &query.StringLiteral{StringValue: "Jane"},
				},
				Op: query.Or,
				Right: &query.FieldExpr{
					Field: "age",
					Op:    query.Equal,
					Value: &query.IntegerLiteral{IntegerValue: 25},
				},
			},
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.expr.Match(target, matcher)
			assert.Equal(t, test.want, result)
		})
	}
}

func TestNotExpr_Match(t *testing.T) {
	target := person{Name: "John", Age: 30}
	matcher := &query.DefaultMatcher{}

	tests := []struct {
		name string
		expr *query.NotExpr
		want bool
	}{
		{
			name: "negate true condition",
			expr: &query.NotExpr{
				Expr: &query.FieldExpr{
					Field: "name",
					Op:    query.Equal,
					Value: &query.StringLiteral{StringValue: "John"},
				},
			},
			want: false,
		},
		{
			name: "negate false condition",
			expr: &query.NotExpr{
				Expr: &query.FieldExpr{
					Field: "name",
					Op:    query.Equal,
					Value: &query.StringLiteral{StringValue: "Jane"},
				},
			},
			want: true,
		},
		{
			name: "negate AND expression",
			expr: &query.NotExpr{
				Expr: &query.BinaryExpr{
					Left: &query.FieldExpr{
						Field: "name",
						Op:    query.Equal,
						Value: &query.StringLiteral{StringValue: "John"},
					},
					Op: query.And,
					Right: &query.FieldExpr{
						Field: "age",
						Op:    query.Equal,
						Value: &query.IntegerLiteral{IntegerValue: 30},
					},
				},
			},
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.expr.Match(target, matcher)
			assert.Equal(t, test.want, result)
		})
	}
}

func TestFieldExpr_Match(t *testing.T) { //nolint:funlen
	target := person{
		Name:     "John",
		Age:      30,
		Height:   1.75,
		IsMember: true,
	}
	matcher := &query.DefaultMatcher{}

	tests := []struct {
		name string
		expr *query.FieldExpr
		want bool
	}{
		{
			name: "string equal - match",
			expr: &query.FieldExpr{
				Field: "name",
				Op:    query.Equal,
				Value: &query.StringLiteral{StringValue: "John"},
			},
			want: true,
		},
		{
			name: "string not equal - match",
			expr: &query.FieldExpr{
				Field: "name",
				Op:    query.NotEqual,
				Value: &query.StringLiteral{StringValue: "Jane"},
			},
			want: true,
		},
		{
			name: "integer greater than - match",
			expr: &query.FieldExpr{
				Field: "age",
				Op:    query.GreaterThan,
				Value: &query.IntegerLiteral{IntegerValue: 25},
			},
			want: true,
		},
		{
			name: "float less than - match",
			expr: &query.FieldExpr{
				Field: "height",
				Op:    query.LessThan,
				Value: &query.NumberLiteral{NumberValue: 1.80},
			},
			want: true,
		},
		{
			name: "non-existent field",
			expr: &query.FieldExpr{
				Field: "invalid",
				Op:    query.Equal,
				Value: &query.StringLiteral{StringValue: "test"},
			},
			want: false,
		},
		{
			name: "field without dumbql tag",
			expr: &query.FieldExpr{
				Field: "IsMember",
				Op:    query.Equal,
				Value: &query.StringLiteral{StringValue: "true"},
			},
			want: false,
		},
		{
			name: "type mismatch",
			expr: &query.FieldExpr{
				Field: "age",
				Op:    query.Equal,
				Value: &query.StringLiteral{StringValue: "30"},
			},
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.expr.Match(target, matcher)
			assert.Equal(t, test.want, result)
		})
	}
}

func TestIdentifier_Match(t *testing.T) { //nolint:funlen
	tests := []struct {
		name   string
		id     query.Identifier
		target any
		op     query.FieldOperator
		want   bool
	}{
		{
			name:   "equal - match",
			id:     query.Identifier("test"),
			target: "test",
			op:     query.Equal,
			want:   true,
		},
		{
			name:   "equal - no match",
			id:     query.Identifier("test"),
			target: "other",
			op:     query.Equal,
			want:   false,
		},
		{
			name:   "not equal - match",
			id:     query.Identifier("test"),
			target: "other",
			op:     query.NotEqual,
			want:   true,
		},
		{
			name:   "not equal - no match",
			id:     query.Identifier("test"),
			target: "test",
			op:     query.NotEqual,
			want:   false,
		},
		{
			name:   "like - match",
			id:     query.Identifier("world"),
			target: "hello world",
			op:     query.Like,
			want:   true,
		},
		{
			name:   "like - no match",
			id:     query.Identifier("universe"),
			target: "hello world",
			op:     query.Like,
			want:   false,
		},
		{
			name:   "with non-string target",
			id:     query.Identifier("42"),
			target: 42,
			op:     query.Equal,
			want:   false,
		},
		{
			name:   "with invalid operator",
			id:     query.Identifier("test"),
			target: "test",
			op:     query.GreaterThan,
			want:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.id.Match(test.target, test.op)
			assert.Equal(t, test.want, result)
		})
	}
}

func TestOneOfExpr_Match(t *testing.T) { //nolint:funlen
	tests := []struct {
		name   string
		expr   *query.OneOfExpr
		target any
		op     query.FieldOperator
		want   bool
	}{
		{
			name: "string equal - match",
			expr: &query.OneOfExpr{
				Values: []query.Valuer{
					&query.StringLiteral{StringValue: "apple"},
					&query.StringLiteral{StringValue: "banana"},
					&query.StringLiteral{StringValue: "orange"},
				},
			},
			target: "banana",
			op:     query.Equal,
			want:   true,
		},
		{
			name: "string equal - no match",
			expr: &query.OneOfExpr{
				Values: []query.Valuer{
					&query.StringLiteral{StringValue: "apple"},
					&query.StringLiteral{StringValue: "banana"},
					&query.StringLiteral{StringValue: "orange"},
				},
			},
			target: "grape",
			op:     query.Equal,
			want:   false,
		},
		{
			name: "integer equal - match",
			expr: &query.OneOfExpr{
				Values: []query.Valuer{
					&query.IntegerLiteral{IntegerValue: 1},
					&query.IntegerLiteral{IntegerValue: 2},
					&query.IntegerLiteral{IntegerValue: 3},
				},
			},
			target: int64(2),
			op:     query.Equal,
			want:   true,
		},
		{
			name: "integer equal - no match",
			expr: &query.OneOfExpr{
				Values: []query.Valuer{
					&query.IntegerLiteral{IntegerValue: 1},
					&query.IntegerLiteral{IntegerValue: 2},
					&query.IntegerLiteral{IntegerValue: 3},
				},
			},
			target: int64(4),
			op:     query.Equal,
			want:   false,
		},
		{
			name: "float equal - match",
			expr: &query.OneOfExpr{
				Values: []query.Valuer{
					&query.NumberLiteral{NumberValue: 1.1},
					&query.NumberLiteral{NumberValue: 2.2},
					&query.NumberLiteral{NumberValue: 3.3},
				},
			},
			target: 2.2,
			op:     query.Equal,
			want:   true,
		},
		{
			name: "float equal - no match",
			expr: &query.OneOfExpr{
				Values: []query.Valuer{
					&query.NumberLiteral{NumberValue: 1.1},
					&query.NumberLiteral{NumberValue: 2.2},
					&query.NumberLiteral{NumberValue: 3.3},
				},
			},
			target: 4.4,
			op:     query.Equal,
			want:   false,
		},
		{
			name: "mixed types",
			expr: &query.OneOfExpr{
				Values: []query.Valuer{
					&query.StringLiteral{StringValue: "one"},
					&query.IntegerLiteral{IntegerValue: 2},
					&query.NumberLiteral{NumberValue: 3.3},
				},
			},
			target: "one",
			op:     query.Equal,
			want:   true,
		},
		{
			name: "empty values",
			expr: &query.OneOfExpr{
				Values: []query.Valuer{},
			},
			target: "test",
			op:     query.Equal,
			want:   false,
		},
		{
			name: "nil values",
			expr: &query.OneOfExpr{
				Values: nil,
			},
			target: "test",
			op:     query.Equal,
			want:   false,
		},
		{
			name: "string like - match",
			expr: &query.OneOfExpr{
				Values: []query.Valuer{
					&query.StringLiteral{StringValue: "world"},
					&query.StringLiteral{StringValue: "universe"},
				},
			},
			target: "hello world",
			op:     query.Like,
			want:   true,
		},
		{
			name: "invalid operator",
			expr: &query.OneOfExpr{
				Values: []query.Valuer{
					&query.StringLiteral{StringValue: "test"},
				},
			},
			target: "test",
			op:     query.GreaterThan,
			want:   false,
		},
		{
			name: "type mismatch",
			expr: &query.OneOfExpr{
				Values: []query.Valuer{
					&query.StringLiteral{StringValue: "42"},
				},
			},
			target: 42,
			op:     query.Equal,
			want:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.expr.Match(test.target, test.op)
			assert.Equal(t, test.want, result)
		})
	}
}
