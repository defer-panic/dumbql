package query

import (
	"reflect"
	"strings"
)

type Matcher interface {
	MatchAnd(target any, left, right Expr) bool
	MatchOr(target any, left, right Expr) bool
	MatchNot(target any, expr Expr) bool
	MatchField(target any, field string, value Valuer, op FieldOperator) bool
	MatchValue(target any, value Valuer, op FieldOperator) bool
}

// DefaultMatcher is a basic implementation of the Matcher interface for evaluating query expressions against structs.
// It supports struct tags using the `dumbql` tag name, which allows you to specify a custom field name.
type DefaultMatcher struct{}

func (m *DefaultMatcher) MatchAnd(target any, left, right Expr) bool {
	return left.Match(target, m) && right.Match(target, m)
}

func (m *DefaultMatcher) MatchOr(target any, left, right Expr) bool {
	return left.Match(target, m) || right.Match(target, m)
}

func (m *DefaultMatcher) MatchNot(target any, expr Expr) bool {
	return !expr.Match(target, m)
}

// MatchField matches a field in the target struct using the provided value and operator. It supports struct tags using
// the `dumbql` tag name, which allows you to specify a custom field name. If struct tag is not provided, it will use
// the field name as is.
func (m *DefaultMatcher) MatchField(target any, field string, value Valuer, op FieldOperator) bool {
	t := reflect.TypeOf(target)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return false
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		fname := f.Name
		if f.Tag.Get("dumbql") != "" {
			fname = f.Tag.Get("dumbql")
		}

		if fname == field {
			v := reflect.ValueOf(target)
			if v.Kind() == reflect.Ptr {
				v = v.Elem()
			}

			return m.MatchValue(v.Field(i).Interface(), value, op)
		}
	}

	return false
}

func (m *DefaultMatcher) MatchValue(target any, value Valuer, op FieldOperator) bool {
	return value.Match(target, op)
}

func (b *BinaryExpr) Match(target any, matcher Matcher) bool {
	switch b.Op {
	case And:
		return matcher.MatchAnd(target, b.Left, b.Right)
	case Or:
		return matcher.MatchOr(target, b.Left, b.Right)
	default:
		return false
	}
}

func (n *NotExpr) Match(target any, matcher Matcher) bool {
	return matcher.MatchNot(target, n.Expr)
}

func (f *FieldExpr) Match(target any, matcher Matcher) bool {
	return matcher.MatchField(target, f.Field.String(), f.Value, f.Op)
}

func (s *StringLiteral) Match(target any, op FieldOperator) bool {
	str, ok := target.(string)
	if !ok {
		return false
	}

	return matchString(str, s.StringValue, op)
}

func (i *IntegerLiteral) Match(target any, op FieldOperator) bool {
	intVal, ok := target.(int64)
	if !ok {
		return false
	}

	return matchNum(intVal, i.IntegerValue, op)
}

func (n *NumberLiteral) Match(target any, op FieldOperator) bool {
	floatVal, ok := target.(float64)
	if !ok {
		return false
	}

	return matchNum(floatVal, n.NumberValue, op)
}

func (i Identifier) Match(target any, op FieldOperator) bool {
	str, ok := target.(string)
	if !ok {
		return false
	}

	return matchString(str, i.String(), op)
}

func (o *OneOfExpr) Match(target any, op FieldOperator) bool {
	switch op { //nolint:exhaustive
	case Equal, Like:
		for _, v := range o.Values {
			if v.Match(target, op) {
				return true
			}
		}

		return false

	default:
		return false
	}
}

func matchString(a, b string, op FieldOperator) bool {
	switch op { //nolint:exhaustive
	case Equal:
		return a == b
	case NotEqual:
		return a != b
	case Like:
		return strings.Contains(a, b)
	default:
		return false
	}
}

func matchNum[T int64 | float64](a, b T, op FieldOperator) bool {
	switch op { //nolint:exhaustive
	case Equal:
		return a == b
	case NotEqual:
		return a != b
	case GreaterThan:
		return a > b
	case GreaterThanOrEqual:
		return a >= b
	case LessThan:
		return a < b
	case LessThanOrEqual:
		return a <= b
	default:
		return false
	}
}
