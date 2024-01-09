# Making an interpreter

This is my first attemp to build an interpreter of my own language, maybe making it 
the succesor for the awfull SL language

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
