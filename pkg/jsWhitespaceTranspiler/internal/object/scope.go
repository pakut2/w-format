package object

type ScopeContext struct {
	For *ForContext
}

type ForContext struct {
	ControlLabelId int64
	EndLabelId     int64
}
