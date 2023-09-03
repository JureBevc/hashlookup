package optional

type Optional[V any] struct {
	Value      V
	IsValueSet bool
}

func Create[T any](a T) *Optional[T] {
	return &Optional[T]{
		Value:      a,
		IsValueSet: true,
	}
}

func Empty[T any]() *Optional[T] {
	var x T
	return &Optional[T]{
		Value:      x,
		IsValueSet: false,
	}
}

type doSomething func()

func (o *Optional[V]) IfSet(fn doSomething) *Optional[V] {
	if o.IsValueSet {
		fn()
	}

	return o
}

func (o *Optional[V]) IfNotSet(fn doSomething) *Optional[V] {
	if !o.IsValueSet {
		fn()
	}

	return o
}
