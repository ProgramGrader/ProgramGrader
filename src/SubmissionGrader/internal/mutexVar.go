package internal

import "sync"

type GoSafeVar[T any] struct {
	mu    sync.Mutex
	value T
}

func GetValue[T any](variable *GoSafeVar[T]) T {
	variable.mu.Lock()
	ret := variable.value
	variable.mu.Unlock()
	return ret
}

func SetValue[T any](variable *GoSafeVar[T], newValue T) {
	variable.mu.Lock()
	variable.value = newValue
	variable.mu.Unlock()
}
