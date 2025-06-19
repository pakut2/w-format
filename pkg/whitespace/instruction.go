package whitespace

import "fmt"

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

func NumberLiteral(value byte) Instruction {
	instruction := Instruction{Body: []Token{SPACE}}

	charBinary := fmt.Sprintf("%s%.8b", instruction, value)

	for _, bit := range charBinary {
		if bit == '1' {
			instruction.Body = append(instruction.Body, TAB)

			continue
		}

		instruction.Body = append(instruction.Body, SPACE)
	}

	instruction.Body = append(instruction.Body, LINE_FEED)

	return instruction
}

func PrintTopStack() Instruction {
	return Instruction{
		Body: []Token{TAB, LINE_FEED, SPACE, SPACE},
	}
}

func EndProgram() Instruction {
	return Instruction{
		Body: []Token{LINE_FEED, LINE_FEED, LINE_FEED},
	}
}
