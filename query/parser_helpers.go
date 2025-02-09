package query

import (
	"fmt"
	"strconv"
)

func resolveBooleanOperator(op any) (BooleanOperator, error) {
	switch string(op.([]byte)) {
	case "AND", "and":
		return And, nil
	case "OR", "or":
		return Or, nil
	default:
		return 0, fmt.Errorf("unknown conditional operator %q", op)
	}
}

func resolveFieldOperator(op any) (FieldOperator, error) {
	switch string(op.([]byte)) {
	case ">=":
		return GreaterThanOrEqual, nil
	case ">":
		return GreaterThan, nil
	case "<=":
		return LessThanOrEqual, nil
	case "<":
		return LessThan, nil
	case "!:", "!=":
		return NotEqual, nil
	case ":", "=":
		return Equal, nil
	default:
		return 0, fmt.Errorf("unknown compare operator %q", op)
	}
}

func resolveOneOfValueType(val any) Valuer {
	switch v := val.(type) {
	case Identifier:
		return &StringLiteral{StringValue: string(v)}
	case string:
		return &StringLiteral{StringValue: v}
	default:
		return v.(Valuer)
	}
}

func parseBooleanExpression(left, rest any) (any, error) {
	expr := left
	for _, r := range rest.([]any) {
		parts := r.([]any)
		// parts[1] holds the operator token, parts[3] holds the next AndExpr.
		//op := string(parts[1].([]byte))
		op, err := resolveBooleanOperator(parts[1])
		if err != nil {
			return nil, err
		}
		right := parts[3]
		expr = &BinaryExpr{
			Left:  expr.(Expr),
			Op:    op,
			Right: right.(Expr),
		}
	}
	return expr, nil
}

func parseFieldExpression(field, op, value any) (any, error) {
	opR, err := resolveFieldOperator(op)
	if err != nil {
		return nil, err
	}

	var val any
	switch v := value.(type) {
	case []byte:
		val = &StringLiteral{StringValue: string(v)}
	case string:
		val = &StringLiteral{StringValue: v}
	case Identifier:
		val = &StringLiteral{StringValue: string(v)}
	default:
		val = value
	}

	return &FieldExpr{
		Field: field.(Identifier),
		Op:    opR,
		Value: val.(Valuer),
	}, nil
}

func parseNumber(c *current) (any, error) {
	if val, err := strconv.Atoi(string(c.text)); err == nil {
		return &IntegerLiteral{IntegerValue: int64(val)}, nil
	}

	if val, err := strconv.ParseFloat(string(c.text), 64); err == nil {
		return &NumberLiteral{NumberValue: val}, nil
	}

	return nil, fmt.Errorf("invalid number literal: %q", string(c.text))
}

func parseString(c *current) (any, error) {
	val, err := strconv.Unquote(string(c.text))
	if err != nil {
		return nil, err
	}
	return &StringLiteral{StringValue: val}, nil
}

func parseOneOfExpression(values any) (any, error) {
	if values == nil || len(values.([]Valuer)) == 0 {
		return &OneOfExpr{Values: nil}, nil
	}

	return &OneOfExpr{Values: values.([]Valuer)}, nil
}

func parseOneOfValues(head, tail any) (any, error) {
	vals := []Valuer{resolveOneOfValueType(head)}

	for _, t := range tail.([]any) {
		// t is an array where index 3 holds the next Value.
		val := resolveOneOfValueType(t.([]any)[3])
		vals = append(vals, val)
	}

	return vals, nil
}
