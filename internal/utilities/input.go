package utilities

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"unicode/utf8"
)

func ReadRune(input *bufio.Reader) rune {
	rune, _, err := input.ReadRune()
	if err != nil {
		if err == io.EOF {
			return 0
		} else {
			panic(fmt.Sprintf("input processing error: %v", err))
		}
	}

	return rune
}

func PeekRune(input bufio.Reader) rune {
	for peekBytes := 4; peekBytes > 0; peekBytes-- {
		peekCharResult, err := input.Peek(peekBytes)
		if err == nil {
			rune, _ := utf8.DecodeRune(peekCharResult)
			if rune == utf8.RuneError {
				return 0
			}

			return rune
		}
	}

	return 0
}

func PeekTwoRunes(input bufio.Reader) (string, error) {
	peekResultBuffer, err := input.Peek(8)
	if err != nil && len(peekResultBuffer) == 0 {
		return "", err
	}

	rune1, rune1Size := utf8.DecodeRune(peekResultBuffer)

	if rune1 == utf8.RuneError && rune1Size == 1 {
		return "", errors.New("malformed input")
	}

	if rune1Size == 0 {
		return "", errors.New("malformed input")
	}

	peekResultChar2Buffer := peekResultBuffer[rune1Size:]

	rune2, rune2Size := utf8.DecodeRune(peekResultChar2Buffer)

	if rune2 == utf8.RuneError && rune2Size == 1 {
		return "", errors.New("malformed input")
	}

	if rune2Size == 0 {
		return "", errors.New("malformed input")
	}

	return string([]rune{rune1, rune2}), nil
}
