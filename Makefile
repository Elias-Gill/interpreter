test:
	gotest ./...

format:
	gofmt -s -w .

repl:
	go run .

lexer:
	go run . --lexer

parser:
	go run . --parser
