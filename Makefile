run-example:
	go run cmd/jsWhitespaceFormatter/main.go -source-file=./examples/source.js -format-file=./examples/format.ts -output-file=./examples/output.ts

test:
	go test ./...
