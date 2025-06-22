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
	sourceFilePath       string
	formatTargetFilePath utilities.Optional[string]
	formatOutputFilePath utilities.Optional[string]
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
	if args.formatTargetFilePath.Valid {
		formatTargetFile, err := os.Open(args.formatTargetFilePath.Value)
		if err != nil {
			panic(fmt.Sprintf("cannot open file: %q, error: %v", args.formatTargetFilePath.Value, err))
		}
		defer formatTargetFile.Close()

		formatTarget = formatTargetFile
	} else {
		formatTarget = strings.NewReader("")
	}

	var formatOutput io.Writer
	if args.formatOutputFilePath.Valid {
		formatOutputFile, err := os.Create(args.formatOutputFilePath.Value)
		if err != nil {
			panic(fmt.Sprintf("cannot open file: %q, error: %v", args.formatOutputFilePath.Value, err))
		}
		defer formatOutputFile.Close()

		formatOutput = formatOutputFile
	} else {
		formatOutput = os.Stdout
	}

	formatter.NewFormatter(formatTarget, whitespace.Instructions(), formatOutput).Format()

	if args.formatOutputFilePath.Valid {
		fmt.Printf("formatted file saved to %q\n", args.formatOutputFilePath.Value)
	}
}

func parseCommandLineArgs() CommandLineArgs {
	sourceFilePath := flag.String("source-file", "", "Whitespace transpilation source file path")
	formatTargetFilePath := flag.String("format-target-file", "", "(Optional) Format target file path. If not provided, outputs generated Whitespace only")
	formatOutputFilePath := flag.String("format-output-file", "", "(Optional) Formatted file output path. If not provided, outputs to stdout")
	flag.Parse()

	if *sourceFilePath == "" {
		panic("source-file not provided")
	}

	parsedFormatTargetFilePath := utilities.Optional[string]{Valid: false, Value: ""}
	if *formatTargetFilePath != "" {
		parsedFormatTargetFilePath = utilities.Optional[string]{
			Valid: true,
			Value: *formatTargetFilePath,
		}
	}

	parsedFormatOutputFilePath := utilities.Optional[string]{Valid: false, Value: ""}
	if *formatOutputFilePath != "" {
		parsedFormatOutputFilePath = utilities.Optional[string]{
			Valid: true,
			Value: *formatOutputFilePath,
		}
	}

	return CommandLineArgs{
		sourceFilePath:       *sourceFilePath,
		formatTargetFilePath: parsedFormatTargetFilePath,
		formatOutputFilePath: parsedFormatOutputFilePath,
	}
}
