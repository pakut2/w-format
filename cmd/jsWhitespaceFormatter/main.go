package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pakut2/w-format/internal/formatter"
	"github.com/pakut2/w-format/internal/utilities"
	"github.com/pakut2/w-format/pkg/jsWhitespaceTranspiler"
)

type CommandLineArgs struct {
	sourceFilePath string
	formatFilePath utilities.Optional[string]
	outputFilePath utilities.Optional[string]
}

func main() {
	args := parseCommandLineArgs()

	sourceFile, err := os.Open(args.sourceFilePath)
	if err != nil {
		panic(fmt.Sprintf("cannot open file: %q, error: %v", args.sourceFilePath, err))
	}
	defer sourceFile.Close()

	lexer := jsWhitespaceTranspiler.NewLexer(sourceFile)
	parsedSource := jsWhitespaceTranspiler.NewParser(lexer).ParseProgram()
	whitespace := jsWhitespaceTranspiler.NewTranspiler().TranspileProgram(parsedSource)

	var formatTarget io.Reader
	if args.formatFilePath.Valid {
		formatTargetFile, err := os.Open(args.formatFilePath.Value)
		if err != nil {
			panic(fmt.Sprintf("cannot open file: %q, error: %v", args.formatFilePath.Value, err))
		}
		defer formatTargetFile.Close()

		formatTarget = formatTargetFile
	} else {
		formatTarget = strings.NewReader("")
	}

	var formatOutput io.Writer
	if args.outputFilePath.Valid {
		formatOutputFile, err := os.Create(args.outputFilePath.Value)
		if err != nil {
			panic(fmt.Sprintf("cannot open file: %q, error: %v", args.outputFilePath.Value, err))
		}
		defer formatOutputFile.Close()

		formatOutput = formatOutputFile
	} else {
		formatOutput = os.Stdout
	}

	formatter.NewFormatter(formatTarget, whitespace.Instructions(), formatOutput).Format()

	if args.outputFilePath.Valid {
		fmt.Printf("output saved to %q\n", args.outputFilePath.Value)
	}
}

func parseCommandLineArgs() CommandLineArgs {
	sourceFilePath := flag.String("source-file", "", "Whitespace transpilation source file path")
	formatFilePath := flag.String("format-file", "", "(Optional) Path to file to be formatted with the generated Whitespace. If not provided, outputs Whitespace only")
	outputFilePath := flag.String("output-file", "", "(Optional) Output file path. If not provided, outputs to stdout")
	flag.Parse()

	if *sourceFilePath == "" {
		panic("source-file not provided")
	}

	parsedFormatTargetFilePath := utilities.Optional[string]{Valid: false, Value: ""}
	if *formatFilePath != "" {
		parsedFormatTargetFilePath = utilities.Optional[string]{
			Valid: true,
			Value: *formatFilePath,
		}
	}

	parsedFormatOutputFilePath := utilities.Optional[string]{Valid: false, Value: ""}
	if *outputFilePath != "" {
		parsedFormatOutputFilePath = utilities.Optional[string]{
			Valid: true,
			Value: *outputFilePath,
		}
	}

	return CommandLineArgs{
		sourceFilePath: *sourceFilePath,
		formatFilePath: parsedFormatTargetFilePath,
		outputFilePath: parsedFormatOutputFilePath,
	}
}
