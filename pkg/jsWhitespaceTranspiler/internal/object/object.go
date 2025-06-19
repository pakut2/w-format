package object

import "github.com/pakut2/w-format/pkg/whitespace"

type ObjectType string

const (
	PROGRAM_OBJ  = "PROGRAM"
	STRING_OBJ   = "STRING"
	CHAR_OBJ     = "CHAR"
	BUILT_IN_OBJ = "BUILT_IN"
	VOID_OBJ     = "VOID"
)

type Object interface {
	Type() ObjectType
	Instructions() []whitespace.Instruction
}

type Program struct {
	WhitespaceInstructions []whitespace.Instruction
}

func (p *Program) Type() ObjectType                       { return PROGRAM_OBJ }
func (p *Program) Instructions() []whitespace.Instruction { return p.WhitespaceInstructions }

type String struct {
	Chars []Char
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Instructions() []whitespace.Instruction {
	return []whitespace.Instruction{}
}

type Char struct {
	HeapAddress byte
}

func (c *Char) Type() ObjectType { return CHAR_OBJ }
func (c *Char) Instructions() []whitespace.Instruction {
	return []whitespace.Instruction{}
}

type BuiltIn struct {
	Function BuiltInFunction
}

func (b *BuiltIn) Type() ObjectType { return BUILT_IN_OBJ }
func (b *BuiltIn) Instructions() []whitespace.Instruction {
	return []whitespace.Instruction{}
}

type BuiltInFunction func(args ...Object) Object

type Void struct{}

func (v *Void) Type() ObjectType { return VOID_OBJ }
func (v *Void) Instructions() []whitespace.Instruction {
	return []whitespace.Instruction{}
}
