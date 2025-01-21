# Language Structure

In this language, every line must end with a semicolon (`;`).

## Variable Declarations

Variables are declared using the reserved word `var`.
They can be used anywhere in the program, provided the scope allows it.

```text
var auxiliar = 2;
var b = 12312;

auxiliar + b;
```

To modify the value of a variable, you must use the `var` keyword again:

```text
var aux = 2;
var aux = aux * 32;

// aux => 64
```

## Comments

The interpreter currently does not support multi-line comments.
Single-line comments are created using `//`.

```text
// This is a comment
var nuevo = 2;
// and this is another comment
```

## If-Else Statements

Conditional statements use the reserved word `si` for "if" and `sino` for "else".
An `elif` statement is not currently available.

```text
si (2 + 2 == 3) {
    return "hola";
} sino {
    return "chau";
}
```

The supported comparison operators are:
- `==` (equal to)
- `<` (less than)
- `>` (greater than)
- `!=` (not equal to)

## Function Declarations, Anonymous Functions, and Function Calls

Functions can be declared as named functions or anonymous functions.
Both use the reserved word `func`.
This allows for higher-order functions.

To return values out of functions, the reserved word `retorna` is used.

```text
var anonima = func() {
    retorna 2;
}

func conNombre(a, b) {
    retorna a() + b;
}

conNombre(anonima, 2);
```

**Note:** Since variable and function identifiers are treated the same, if you declare a
function and then a variable with the same name, the variable will overwrite the function.
For example:

```text
func bar() {...}
var bar = 2;
// bar() can no longer be used.
```

Functions support recursion and can be passed as parameters to other functions.

```text
func Fibonacci(n) {
    si (n < 1) {
        retorna n;
    }
    si (n == 1) {
        retorna n;
    }
    retorna Fibonacci(n-1) + Fibonacci(n-2);
}

var bar = func() { retorna Fibonacci(8); };

var baz = func (a) {retorna a();};

baz(bar);
```

# Making an Interpreter

This is my first attempt at building an interpreter.
This doesn't aim for complex or innovative features, nor performance or any kind of
optimization.
Just for learning purposes about parsing and evaluation of text.

## Running and Testing

A simple **REPL** (Read-Eval-Print Loop) is provided for interactive use.
To run it:

```text
go run .
```

This REPL has 3 modes.
To see the available modes run `go run . -help`.

To run the tests suit, use the standard Go test command:

```text
go test ./...
```
