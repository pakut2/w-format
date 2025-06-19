package formatter

import (
	"fmt"
	"strings"

	"github.com/pakut2/js-whitespace/pkg/whitespace"
)

type Formatter struct {
	input                            string
	whitespaceInstructionTokens      []whitespace.Token
	whitespaceFinalInstructionTokens []whitespace.Token
	whitespaceTokenIndex             int
	position                         int
	readPosition                     int
	currentChar                      byte
}

func NewFormatter(input string, whitespaceInstructions []whitespace.Instruction) *Formatter {
	f := &Formatter{input: input, whitespaceTokenIndex: 0}

	whitespaceInstructionsLength := len(whitespaceInstructions)

	for _, instruction := range whitespaceInstructions[:whitespaceInstructionsLength-1] {
		f.whitespaceInstructionTokens = append(f.whitespaceInstructionTokens, instruction.Body...)
	}

	f.whitespaceFinalInstructionTokens = whitespaceInstructions[whitespaceInstructionsLength-1].Body

	f.readChar()

	return f
}

func (f *Formatter) getNextWhitespaceToken() whitespace.Token {
	if f.whitespaceTokenIndex >= len(f.whitespaceInstructionTokens)-1 {
		f.whitespaceInstructionTokens = append(f.whitespaceInstructionTokens, whitespace.Noop().Body...)
	}

	whitespaceToken := f.whitespaceInstructionTokens[f.whitespaceTokenIndex]

	f.whitespaceTokenIndex++

	return whitespaceToken
}

func (f *Formatter) getNextWhitespaceTokenUntil(target whitespace.Token) []whitespace.Token {
	initialWhitespaceTokenIndex := f.whitespaceTokenIndex

	for _, whitespaceToken := range f.whitespaceInstructionTokens[initialWhitespaceTokenIndex:] {
		if f.whitespaceTokenIndex >= len(f.whitespaceInstructionTokens)-1 {
			break
		}

		f.whitespaceTokenIndex++

		if whitespaceToken == target {
			break
		}
	}

	return f.whitespaceInstructionTokens[initialWhitespaceTokenIndex:f.whitespaceTokenIndex]
}

func (f *Formatter) readChar() {
	if f.readPosition >= len(f.input) {
		f.currentChar = 0
	} else {
		f.currentChar = f.input[f.readPosition]
	}

	f.position = f.readPosition
	f.readPosition += 1
}

func (f *Formatter) Format() string {
	var formattedInput string

	for f.currentChar != 0 {
		switch f.currentChar {
		case ' ', '\t':
			formattedInput = fmt.Sprintf("%s%c", formattedInput, f.getNextWhitespaceToken())
		case '\n':
			formattedInput = fmt.Sprintf("%s%s", formattedInput, string(f.getNextWhitespaceTokenUntil(whitespace.LINE_FEED)))
		case '"', '\'', '`':
			text := f.readString()

			formattedInput = fmt.Sprintf("%s%c%s%c", formattedInput, f.currentChar, strings.ReplaceAll(text, " ", "\u2007"), f.currentChar)
		default:
			formattedInput = fmt.Sprintf("%s%c", formattedInput, f.currentChar)
		}

		f.readChar()
	}

	if f.whitespaceTokenIndex < len(f.whitespaceInstructionTokens) {
		formattedInput = fmt.Sprintf(
			"%s%s",
			formattedInput,
			string(f.whitespaceInstructionTokens[f.whitespaceTokenIndex:len(f.whitespaceInstructionTokens)]),
		)
	}

	formattedInput = fmt.Sprintf("%s%s", formattedInput, string(f.whitespaceFinalInstructionTokens))

	return formattedInput
}

// TODO don't break on different quote than string start
func (f *Formatter) readString() string {
	position := f.position + 1

	for {
		f.readChar()
		if f.currentChar == '"' || f.currentChar == '\'' || f.currentChar == '`' || f.currentChar == 0 {
			break
		}
	}

	return f.input[position:f.position]
}
