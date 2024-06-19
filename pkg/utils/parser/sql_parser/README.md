## The SQL Parser

The SQL parser parses and validates a SQL string.
**WARNING** This version of the code does not pretend to be a complete SQL parser. It is currently intended to parse only WHERE clauses.

It parses the string by feeding a SQL grammar and a SQLScanner to the `StringParser` object.

Additionally, it will return two values that you can use to pass the SQL string to your database. 

Those values are:
* Query string: this is the same as the received query, but all the values are replaced with `?`, so that you can feed the prepared statement to the DB
* Values []interface{}: this contains all the values to be passed to the DB, in the right order , to replace the `?`

For example, parsing the following SQL string

```sql
COMPANY_NAME='Red Hat' and COUNTRY='Ireland'
```
you will get:
```sql
Query: "COMPANY_NAME = ? and COUNTRY = ?
Values: "Red Hat", "Ireland"
```

### Instantiating the parser
The parser uses the `functional options` pattern. Instantiating it with all the defaults is as easy as calling one function:
```go
parser := NewSQLParser()
```

The `NewSQLParser` function takes a variadic list of `SQLParserOption` that can be passed to configure the parser instance.

#### Supported options
##### WithValidColumns( validColumns ...string)
This can be used to limit the column the user can insert into the SQL string.
For example, this will lead to a validation error
```go
parser := NewSQLParser(WithValidColumns("surname"))
_, _, err := parser.Parse("name = 'mickey' and surname = 'mouse'")
fmt.Println(err)

---- output

[1] error parsing the filter: invalid column name: 'name', valid values are: [surname]
```
The number in the square bracket represent the position in the string where the error occurred.

##### WithMaximumComplexity( maximumComplexity int )
This can be used to specify the maximum number of logical operator allowed into the query
```go
parser := NewSQLParser(
    WithMaximumComplexity(2),
)
_, _, err := parser.Parse("(name = 'mickey' or name = 'minnie') and surname = 'mouse' and age > 20")
fmt.Println(err)

---- output

[60] error parsing the filter: maximum number of permitted joins (2) exceeded
```
##### WithColumnPrefix(columnPrefix string)
This option specifies the prefix to be added to each column in the produced output qry.
For example, if we want every column to be prefixed with 'main.', we will use the following code
```go
	parser := NewSQLParser(WithColumnPrefix("main"))
qry, _, _ := parser.Parse("(name = 'mickey' or name = 'minnie') and surname = 'mouse' and age >= 20")
fmt.Println(qry)

---- output

(main.name = ? or main.name = ?) and main.surname = ? and main.age >= ?
```
##### All the options together
```go
parser := NewSQLParser(
    WithValidColumns("surname"),
    WithColumnPrefix("main"),
    WithMaximumComplexity(2),
)
qry, _, err := parser.Parse("(name = 'mickey' or name = 'minnie') and surname = 'mouse' and age >= 20")
fmt.Println("err: ", err)
fmt.Println("qry: ", qry)

---- output

err:  [2] error parsing the filter: invalid column name: 'name', valid values are: [surname age]
qry:
```
