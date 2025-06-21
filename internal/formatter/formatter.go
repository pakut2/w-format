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
	input bufio.Reader

	whitespaceInstructionTokens      []whitespace.Token
	whitespaceFinalInstructionTokens []whitespace.Token
	whitespaceTokenIndex             int

	previousChar rune
	currentChar  rune
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
	f.previousChar = f.currentChar
	f.currentChar = utilities.ReadRune(&f.input)
}

func (f *Formatter) Format(target io.Writer) {
	formattedOutput := bufio.NewWriter(target)

	for f.currentChar != 0 {
		switch f.currentChar {
		case ' ', '\t':
			nextTwoChars, err := utilities.PeekTwoRunes(f.input)
			if err == nil && nextTwoChars == "=>" {
				if f.peekNextWhitespaceToken() == whitespace.LINE_FEED {
					if _, err = formattedOutput.WriteString("\u2007"); err != nil {
						f.handleOutputError(err)
					}
				}
			} else {
				if err := formattedOutput.WriteByte(byte(f.getNextWhitespaceToken())); err != nil {
					f.handleOutputError(err)
				}
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
		case '/':
			if utilities.PeekRune(f.input) == '/' {
				f.readChar()
				commentLiteral := f.readComment()

				_, err := formattedOutput.WriteString(fmt.Sprintf("//%s%s", f.sanitizeString(commentLiteral), f.getNextWhitespaceTokenUntil(whitespace.LINE_FEED)))
				if err != nil {
					f.handleOutputError(err)
				}
			} else if utilities.PeekRune(f.input) == '*' {
				f.readChar()
				blockCommentLiteral := f.readBlockComment()

				_, err := formattedOutput.WriteString(fmt.Sprintf("/*%s*/", f.sanitizeString(blockCommentLiteral)))
				if err != nil {
					f.handleOutputError(err)
				}

				f.readChar()
			} else {
				_, err := formattedOutput.WriteRune(f.currentChar)
				if err != nil {
					f.handleOutputError(err)
				}
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

	if err = formattedOutput.Flush(); err != nil {
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

func (f *Formatter) handleOutputError(err error) {
	panic(fmt.Sprintf("cannot write formatted output, error: %v", err))
}
