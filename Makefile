run:
	go run cmd/jsWhitespaceFormatter/main.go -source-file=./examples/source.js -format-target-file=./examples/formatTarget.ts -format-output-file=./examples/formatOutput.ts

test:
	go test ./...
