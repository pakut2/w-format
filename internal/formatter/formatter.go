package formatter

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/pakut2/js-whitespace/pkg/whitespace"
)

type Formatter struct {
	input                            bufio.Reader
	whitespaceInstructionTokens      []whitespace.Token
	whitespaceFinalInstructionTokens []whitespace.Token
	whitespaceTokenIndex             int
	currentChar                      rune
}

func NewFormatter(input io.Reader, whitespaceInstructions []whitespace.Instruction) *Formatter {
	f := &Formatter{input: *bufio.NewReader(input), whitespaceTokenIndex: 0}

	whitespaceInstructionsLength := len(whitespaceInstructions)

	for _, instruction := range whitespaceInstructions[:whitespaceInstructionsLength-1] {
		f.whitespaceInstructionTokens = append(f.whitespaceInstructionTokens, instruction.Body...)
	}

	f.whitespaceFinalInstructionTokens = whitespaceInstructions[whitespaceInstructionsLength-1].Body

	f.readChar()

	return f
}

func (f *Formatter) readChar() {
	char, _, err := f.input.ReadRune()
	if err != nil {
		if err == io.EOF {
			f.currentChar = 0

			return
		} else {
			panic(err)
		}
	}

	f.currentChar = char
}

func (f *Formatter) Format(target io.Writer) {
	formattedOutput := bufio.NewWriter(target)

	for f.currentChar != 0 {
		switch f.currentChar {
		case ' ', '\t':
			err := formattedOutput.WriteByte(byte(f.getNextWhitespaceToken()))
			if err != nil {
				f.handleOutputError(err)
			}
		case '\n':
			_, err := formattedOutput.WriteString(string(f.getNextWhitespaceTokenUntil(whitespace.LINE_FEED)))
			if err != nil {
				f.handleOutputError(err)
			}
		case '"', '\'', '`':
			stringLiteral := f.readString()

			_, err := formattedOutput.WriteString(fmt.Sprintf("%c%s%c", f.currentChar, f.sanitizeString(stringLiteral), f.currentChar))
			if err != nil {
				f.handleOutputError(err)
			}
		default:
			_, err := formattedOutput.WriteRune(f.currentChar)
			if err != nil {
				f.handleOutputError(err)
			}
		}

		f.readChar()
	}

	if f.whitespaceTokenIndex < len(f.whitespaceInstructionTokens) {
		_, err := formattedOutput.WriteString(string(f.whitespaceInstructionTokens[f.whitespaceTokenIndex:len(f.whitespaceInstructionTokens)]))
		if err != nil {
			f.handleOutputError(err)
		}
	}

	_, err := formattedOutput.WriteString(string(f.whitespaceFinalInstructionTokens))
	if err != nil {
		f.handleOutputError(err)
	}

	err = formattedOutput.Flush()
	if err != nil {
		f.handleOutputError(err)
	}
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

// TODO don't break on different quote than string start
func (f *Formatter) readString() string {
	var stringLiteral string

	for {
		f.readChar()

		if f.currentChar == '"' || f.currentChar == '\'' || f.currentChar == '`' || f.currentChar == 0 {
			break
		}

		stringLiteral = fmt.Sprintf("%s%c", stringLiteral, f.currentChar)
	}

	return stringLiteral
}

func (f *Formatter) sanitizeString(value string) string {
	sanitizedString := strings.ReplaceAll(value, " ", "\u2007")
	sanitizedString = strings.ReplaceAll(sanitizedString, "\t", strings.Repeat("\u2007", 4))
	sanitizedString = strings.ReplaceAll(sanitizedString, "\n", "\u000a")

	return sanitizedString
}

func (f *Formatter) handleOutputError(err error) {
	panic(fmt.Sprintf("cannot write formatted output, error: %v", err))
}
