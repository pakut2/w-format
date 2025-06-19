package utilities

type Optional[T any] struct {
	Valid bool
	Value T
}
