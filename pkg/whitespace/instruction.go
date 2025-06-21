package whitespace

import (
	"strconv"
)

type Token byte

const (
	TAB       = '\t'
	LINE_FEED = '\n'
	SPACE     = ' '
	//TAB       = 'T'
	//LINE_FEED = 'L'
	//SPACE     = 'S'
)

type Instruction struct {
	Body []Token
}

func (i *Instruction) String() string {
	return string(i.Body)
}

func Noop() Instruction {
	return Instruction{
		Body: []Token{SPACE, TAB, LINE_FEED, SPACE, LINE_FEED},
	}
}

func StoreInHeap() Instruction {
	return Instruction{
		Body: []Token{TAB, TAB, SPACE},
	}
}

func RetrieveFromHeap() Instruction {
	return Instruction{
		Body: []Token{TAB, TAB, TAB},
	}
}

func PushToStack() Instruction {
	return Instruction{
		Body: []Token{SPACE, SPACE},
	}
}

func SwapTwoTopStackItems() Instruction {
	return Instruction{
		Body: []Token{SPACE, LINE_FEED, TAB},
	}
}

func LiftStackItem(itemOrdinalNumber int) Instruction {
	itemOrdinalNumberLiteral := NumberLiteral(int64(itemOrdinalNumber))

	return Instruction{
		Body: append([]Token{SPACE, TAB, SPACE}, itemOrdinalNumberLiteral.Body...),
	}
}

func NumberLiteral(value int64) Instruction {
	var instruction Instruction

	if value >= 0 {
		instruction = Instruction{Body: []Token{SPACE}}
	} else {
		instruction = Instruction{Body: []Token{TAB}}
	}

	binaryNumber := strconv.FormatInt(value, 2)

	for _, bit := range binaryNumber {
		if bit == '1' {
			instruction.Body = append(instruction.Body, TAB)

			continue
		}

		instruction.Body = append(instruction.Body, SPACE)
	}

	instruction.Body = append(instruction.Body, LINE_FEED)

	return instruction
}

func PrintTopStackChar() Instruction {
	return Instruction{
		Body: []Token{TAB, LINE_FEED, SPACE, SPACE},
	}
}

func PrintTopStackInteger() Instruction {
	return Instruction{
		Body: []Token{TAB, LINE_FEED, SPACE, TAB},
	}
}

func Add() Instruction {
	return Instruction{
		Body: []Token{TAB, SPACE, SPACE, SPACE},
	}
}

func Subtract() Instruction {
	return Instruction{
		Body: []Token{TAB, SPACE, SPACE, TAB},
	}
}

func Multiply() Instruction {
	return Instruction{
		Body: []Token{TAB, SPACE, SPACE, LINE_FEED},
	}
}

func Divide() Instruction {
	return Instruction{
		Body: []Token{TAB, SPACE, TAB, SPACE},
	}
}

func Modulo() Instruction {
	return Instruction{
		Body: []Token{TAB, SPACE, TAB, TAB},
	}
}

func Label(labelId int64) Instruction {
	labelIdLiteral := NumberLiteral(labelId)

	return Instruction{
		Body: append([]Token{LINE_FEED, SPACE, SPACE}, labelIdLiteral.Body...),
	}
}

func JumpToLabel(labelId int64) Instruction {
	labelIdLiteral := NumberLiteral(labelId)

	return Instruction{
		Body: append([]Token{LINE_FEED, SPACE, LINE_FEED}, labelIdLiteral.Body...),
	}
}

func JumpToLabelIfZero(labelId int64) Instruction {
	labelIdLiteral := NumberLiteral(labelId)

	return Instruction{
		Body: append([]Token{LINE_FEED, TAB, SPACE}, labelIdLiteral.Body...),
	}
}

func JumpToLabelIfNegative(labelId int64) Instruction {
	labelIdLiteral := NumberLiteral(labelId)

	return Instruction{
		Body: append([]Token{LINE_FEED, TAB, TAB}, labelIdLiteral.Body...),
	}
}

func EndProgram() Instruction {
	return Instruction{
		Body: []Token{LINE_FEED, LINE_FEED, LINE_FEED},
	}
}
