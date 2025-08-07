package pool

type Query[T comparable] struct {
	Results []any
	Match   func(item any) bool
}
