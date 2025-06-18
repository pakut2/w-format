package whitespace

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
