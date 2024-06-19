## String Scanner

A string scanner is a lightweight parsing tool that iterates over a string, segmenting it into discrete 
units based on a defined delimiter pattern. The most basic implementation treats each character as a 
distinct token.

A String Scanner must implement the interface below:

```go
type Scanner interface {
	// Next - Move to the next Token. Return false if no next Token is available
	Next() bool
	// Peek - Look at the next Token without moving. Return false if no next Token is available
	Peek() (bool, *Token)
	// Token - Return the current Token Value. Panics if current Position is invalid.
	Token() *Token
	// Init - Initialise the scanner with the given string
	Init(s string)
}
```

This package provides two implementation:
* SimpleStringScanner: this is the simplest implementation. It just iterates over each character of the provided string
* SQLStringScanner: this scanner splits the string into tokens that can be used to parse the string as a SQL string.

### Example usage

```go
scanner := string_scanner.NewSimpleScanner()
scanner.Init("SELECT * FROM ADDRESS_BOOK WHERE COMPANY='RED HAT'")
for scanner.Next() {
	fmt.Println(scanner.Token().Value)
}
```

The code above prints all the tokens:
```go
S
E
L
E
C
T
 
*
 
F
...cut...
```
Using the SQLScanner each token will be one SQL element:
```go
scanner := sql_parser.NewSQLScanner()
scanner.Init("SELECT * FROM ADDRESS_BOOK WHERE COMPANY='RED HAT'")
for scanner.Next() {
	fmt.Println(scanner.Token().Value)
}
```
output:
```
SELECT
*
FROM
ADDRESS_BOOK
WHERE
COMPANY
=
'RED HAT'
```
