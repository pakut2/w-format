package object

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
	Instruction() string
}

type Program struct {
	InstructionBody string
}

func (p *Program) Type() ObjectType    { return PROGRAM_OBJ }
func (p *Program) Instruction() string { return p.InstructionBody }

type String struct {
	InstructionBody string
	Chars           []Char
}

func (s *String) Type() ObjectType    { return STRING_OBJ }
func (s *String) Instruction() string { return s.InstructionBody }

type Char struct {
	HeapAddress byte
}

func (c *Char) Type() ObjectType    { return CHAR_OBJ }
func (c *Char) Instruction() string { return "" }

type BuildIn struct {
	Function BuiltInFunction
}

func (b *BuildIn) Type() ObjectType    { return BUILT_IN_OBJ }
func (b *BuildIn) Instruction() string { return "" }

type BuiltInFunction func(args ...Object) Object

type Void struct {
	InstructionBody string
}

func (v *Void) Type() ObjectType    { return VOID_OBJ }
func (v *Void) Instruction() string { return v.InstructionBody }
