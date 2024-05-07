# Language structure

This language is made for mimic SL2 syntax but in a C-like style.
So after every line a colon (;) is necesary.

#### Variable declarations
Variables can be created with the reserved word "var", and can be used everywhere on the
program (as long as the scope allows it)
```text
var auxiliar = 2;
var b = 12312;

auxiliar + b
```
#### If-else statements

Like SL2, if statements use the reserverd word "si" and for else statements uses "sino".
An else statement with a condition is currently unavailable.

```text
si (2 + 2 == 3) {
    return "hola";
} sino {
    return "chau";
}
```

#### Function declarations, anonymous functions and function calls

Now are currently supported function declarations as statements or as anonymous functions, so
we can have higher order functions.
Both are created with the reserved word "func".

```
var anonima = func() {
    return 2;
}

func conNombre(a, b) {
    return a() + b;
}

conNombre(anonima, 2)
```

# Making an interpreter

This is my first attemp to build an interpreter of my own language, maybe making it the
succesor for the awfull SL language

## Util dependencies
This project relies enterily from the golang standard library, but, to enhance the development
experience I recommend:
- GNU Make tool.
- [gotest](https://github.com/rakyll/gotest) for colored test output.

## Running and testing
I probide a simple _REPL_ which can be run with:
```
go run .
```

For testing:
```
make test       // if you have make and gotest installed
go test ./...   // default go test command
```
