// Provides a set of tools and syntax sugar context.
package gocontext

import (
	"context"
	"sync"

	"golang.org/x/exp/constraints"
)

type (
	goContextKey string
	contextValue any
)

const (
	contextKey goContextKey = "gocontext-values"
)

// FromContext Creates a new child context
func FromContext(ctx context.Context) context.Context {
	if _, ok := ctx.Value(contextKey).(*sync.Map); ok {
		return ctx
	}
	return context.WithValue(ctx, contextKey, &sync.Map{})
}

// Add Add a new value into context
func Add[T constraints.Integer | constraints.Float | ~string | ~bool](ctx context.Context, key string, value T) context.Context {
	if syncMap, ok := ctx.Value(contextKey).(*sync.Map); ok {
		syncMap.Store(key, value)
		return ctx
	}
	return Add(FromContext(ctx), key, value)
}

// Get Get a value from context
func Get[T constraints.Integer | constraints.Float | ~string | ~bool](ctx context.Context, key string) (T, bool) {
	if syncMap, ok := ctx.Value(contextKey).(*sync.Map); ok {
		if rawValue, ok := syncMap.Load(key); ok {
			ctxValue := any(rawValue).(contextValue)
			return any(ctxValue).(T), true
		}
	}
	var result T
	return result, false
}

// GetKeys Get all context keys
func GetKeys(ctx context.Context) (keys []string) {
	if syncMap, ok := ctx.Value(contextKey).(*sync.Map); ok {
		syncMap.Range(func(key, value any) bool {
			keys = append(keys, key.(string))
			return true
		})
	}
	return
}

// Delete Remove a context value
func Delete(ctx context.Context, key string) {
	if syncMap, ok := ctx.Value(contextKey).(*sync.Map); ok {
		syncMap.Delete(key)
	}
}
