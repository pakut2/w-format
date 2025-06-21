package formatter

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/pakut2/w-format/internal/utilities"
	"github.com/pakut2/w-format/pkg/whitespace"
)

type Formatter struct {
	input  bufio.Reader
	target bufio.Writer

	whitespaceInstructionTokens      []whitespace.Token
	whitespaceFinalInstructionTokens []whitespace.Token
	whitespaceTokenIndex             int

	previousChar rune
	currentChar  rune
}

func NewFormatter(input io.Reader, whitespaceInstructions []whitespace.Instruction, target io.Writer) *Formatter {
	f := &Formatter{input: *bufio.NewReader(input), target: *bufio.NewWriter(target), whitespaceTokenIndex: 0}

	whitespaceInstructionsLength := len(whitespaceInstructions)

	for _, instruction := range whitespaceInstructions[:whitespaceInstructionsLength-1] {
		f.whitespaceInstructionTokens = append(f.whitespaceInstructionTokens, instruction.Body...)
	}

	f.whitespaceFinalInstructionTokens = whitespaceInstructions[whitespaceInstructionsLength-1].Body

	f.readChar()

	return f
}

func (f *Formatter) readChar() {
	f.previousChar = f.currentChar
	f.currentChar = utilities.ReadRune(&f.input)
}

func (f *Formatter) Format() {
	for f.currentChar != 0 {
		switch f.currentChar {
		case ' ', '\t':
			nextTwoChars, err := utilities.PeekTwoRunes(f.input)
			if err == nil && nextTwoChars == "=>" {
				if f.peekNextWhitespaceToken() == whitespace.LINE_FEED {
					f.writeString("\u2007")
				}
			} else {
				if err = f.target.WriteByte(byte(f.getNextWhitespaceToken())); err != nil {
					f.handleOutputError(err)
				}
			}
		case '\n':
			f.writeString(string(f.getNextWhitespaceTokenUntil(whitespace.LINE_FEED)))
		case '"', '\'', '`':
			stringLiteral := f.readString()

			f.writeString(fmt.Sprintf("%c%s%c", f.currentChar, f.sanitizeString(stringLiteral), f.currentChar))
		case '/':
			if utilities.PeekRune(f.input) == '/' {
				f.readChar()
				commentLiteral := f.readComment()

				f.writeString(
					fmt.Sprintf(
						"//%s%s",
						f.sanitizeString(commentLiteral),
						f.getNextWhitespaceTokenUntil(whitespace.LINE_FEED),
					),
				)
			} else if utilities.PeekRune(f.input) == '*' {
				f.readChar()
				blockCommentLiteral := f.readBlockComment()

				f.writeString(fmt.Sprintf("/*%s*/", f.sanitizeString(blockCommentLiteral)))

				f.readChar()
			} else {
				f.writeChar(f.currentChar)
			}
		default:
			f.writeChar(f.currentChar)
		}

		f.readChar()
	}

	if f.whitespaceTokenIndex < len(f.whitespaceInstructionTokens) {
		f.writeString(string(f.whitespaceInstructionTokens[f.whitespaceTokenIndex:len(f.whitespaceInstructionTokens)]))
	}

	f.writeString(string(f.whitespaceFinalInstructionTokens))

	if err := f.target.Flush(); err != nil {
		f.handleOutputError(err)
	}
}

func (f *Formatter) peekNextWhitespaceToken() whitespace.Token {
	nextTokenIndex := f.whitespaceTokenIndex + 1

	if nextTokenIndex >= len(f.whitespaceInstructionTokens)-1 {
		return whitespace.Noop().Body[0]
	}

	return f.whitespaceInstructionTokens[nextTokenIndex]
}

func (f *Formatter) getNextWhitespaceToken() whitespace.Token {
	if f.whitespaceTokenIndex >= len(f.whitespaceInstructionTokens)-1 {
		f.whitespaceInstructionTokens = append(f.whitespaceInstructionTokens, whitespace.Noop().Body...)
	}

	token := f.whitespaceInstructionTokens[f.whitespaceTokenIndex]

	f.whitespaceTokenIndex++

	return token
}

func (f *Formatter) getNextWhitespaceTokenUntil(target whitespace.Token) []whitespace.Token {
	var tokens []whitespace.Token

	for {
		currentToken := f.getNextWhitespaceToken()
		tokens = append(tokens, currentToken)

		if currentToken == target {
			break
		}
	}

	return tokens
}

func (f *Formatter) readString() string {
	startingQuote := f.currentChar

	var stringLiteral string

	for {
		f.readChar()

		if f.currentChar == 0 || (f.currentChar == startingQuote && f.previousChar != '\\') {
			break
		}

		stringLiteral = fmt.Sprintf("%s%c", stringLiteral, f.currentChar)
	}

	return stringLiteral
}

func (f *Formatter) readComment() string {
	var commentLiteral string

	for {
		f.readChar()

		if f.currentChar == 0 || f.currentChar == '\n' {
			break
		}

		commentLiteral = fmt.Sprintf("%s%c", commentLiteral, f.currentChar)
	}

	return commentLiteral
}

func (f *Formatter) readBlockComment() string {
	var blockCommentLiteral string

	for {
		f.readChar()

		if f.currentChar == 0 || (f.currentChar == '*' && utilities.PeekRune(f.input) == '/') {
			break
		}

		blockCommentLiteral = fmt.Sprintf("%s%c", blockCommentLiteral, f.currentChar)
	}

	return blockCommentLiteral
}

func (f *Formatter) sanitizeString(value string) string {
	sanitizedString := strings.ReplaceAll(value, " ", "\u2007")
	sanitizedString = strings.ReplaceAll(sanitizedString, "\t", strings.Repeat("\u2007", 4))
	sanitizedString = strings.ReplaceAll(sanitizedString, "\n", "\u2028")

	return sanitizedString
}

func (f *Formatter) writeChar(char rune) {
	_, err := f.target.WriteRune(char)
	if err != nil {
		f.handleOutputError(err)
	}
}

func (f *Formatter) writeString(value string) {
	_, err := f.target.WriteString(value)
	if err != nil {
		f.handleOutputError(err)
	}
}

func (f *Formatter) handleOutputError(err error) {
	panic(fmt.Sprintf("cannot write formatted output, error: %v", err))
}
