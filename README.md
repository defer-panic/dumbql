# dumbql

Simple (dumb) query language and parser for Go. 

## Table of Contents

<!-- TOC -->
* [dumbql](#dumbql)
  * [Table of Contents](#table-of-contents)
  * [Features](#features)
  * [Examples](#examples)
    * [Simple parse](#simple-parse)
    * [Validation against schema](#validation-against-schema)
    * [Convert to SQL](#convert-to-sql)
  * [Query syntax](#query-syntax)
    * [Field expression](#field-expression)
    * [Field expression operators](#field-expression-operators)
    * [Boolean operators](#boolean-operators)
    * [“One of” expression](#one-of-expression)
    * [Numbers](#numbers)
    * [Strings](#strings)
<!-- TOC -->

## Features

- Field expressions (`age >= 18`, `field.name:"field value"`, etc.)
- Boolean expressions (`age >= 18 and city = Barcelona`, `occupation = designer or occupation = "ux analyst"`)
- One-of/In expressions (`occupation = [designer, "ux analyst"]`)
- Schema validation
- Drop-in usage with [squirrel](https://github.com/Masterminds/squirrel) query builder or SQL drivers directly

## Examples

### Simple parse

```go
package main

import (
    "fmt"

    "github.com/defer-panic/dumbql"
)

func main() {
    const q = `profile.age >= 18 and profile.city = Barcelona`
    ast, err := dumbql.Parse(q)
    if err != nil {
        panic(err)
    }

    fmt.Println(ast)
    // Output: (and (>= profile.age 18) (= profile.city "Barcelona"))
}
```

### Validation against schema

```go
package main

import (
    "fmt"

    "github.com/defer-panic/dumbql"
    "github.com/defer-panic/dumbql/schema"
)

func main() {
    schm := schema.Schema{
        "status": schema.All(
            schema.Is[string](),
            schema.EqualsOneOf("pending", "approved", "rejected"),
        ),
        "period_months": schema.Max(int64(3)),
        "title":         schema.LenInRange(1, 100),
    }

    // The following query is invalid against the schema:
    // 	- period_months == 4, but max allowed value is 3
    // 	- field `name` is not described in the schema
    //
    // Invalid parts of the query are dropped.
    const q = `status:pending and period_months:4 and (title:"hello world" or name:"John Doe")`
    expr, err := dumbql.Parse(q)
    if err != nil {
        panic(err)
    }

    validated, err := expr.Validate(schm)
    fmt.Println(validated)
    fmt.Printf("validation error: %v\n", err)
    // Output: 
    // (and (= status "pending") (= title "hello world"))
    // validation error: field "period_months": value must be equal or less than 3, got 4; field "name" not found in schema
}
```

### Convert to SQL

```go
package main

import (
  "fmt"

  sq "github.com/Masterminds/squirrel"
  "github.com/defer-panic/dumbql"
)

func main() {
  const q = `status:pending and period_months < 4 and (title:"hello world" or name:"John Doe")`
  expr, err := dumbql.Parse(q)
  if err != nil {
    panic(err)
  }

  sql, args, err := sq.Select("*").
    From("users").
    Where(expr).
    ToSql()
  if err != nil {
    panic(err)
  }

  fmt.Println(sql)
  fmt.Println(args)
  // Output: 
  // SELECT * FROM users WHERE ((status = ? AND period_months < ?) AND (title = ? OR name = ?))
  // [pending 4 hello world John Doe]
}

```

See [dumbql_example_test.go](dumbql_example_test.go)

## Query syntax

This section is a non-formal description of DumbQL syntax. For strict description see [grammar file](query/grammar.peg).

### Field expression

Field name & value pair divided by operator. Field name is any alphanumeric identifier (with underscore), value can be string, int64 or floa64.
One-of expression is also supported (see below).

```
<field_name> <operator> <value>
```

for example

```
period_months < 4
```

### Field expression operators

| Operator             | Meaning       | Supported types              |
|----------------------|---------------|------------------------------|
| `:` or `=`           | Equal, one of | `int64`, `float64`, `string` |
| `!=` or `!:`         | Not equal     | `int64`, `float64`, `string` |
| `>`, `>=`, `<`, `<=` | Comparison    | `int64`, `float64`           |


### Boolean operators

Multiple field expression can be combined into boolean expressions with `and` (`AND`) or `or` (`OR`) operators:

```
status:pending and period_months < 4 and (title:"hello world" or name:"John Doe")
```

### “One of” expression

Sometimes instead of multiple `and`/`or` clauses against the same field:

```
occupation = designer or occupation = "ux analyst"
```

it's more convenient to use equivalent “one of” expressions:

```
occupation: [designer, "ux analyst"]
```

### Numbers

If number does not have digits after `.` it's treated as integer and stored as `int64`. And it's `float64` otherwise.

### Strings

String is a sequence on Unicode characters surrounded by double quotes (`"`). In some cases like single word it's possible to write string value without double quotes.
