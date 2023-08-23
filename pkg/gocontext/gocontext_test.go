package gocontext

import (
	"context"
	"sync"
	"testing"

	"golang.org/x/exp/slices"
)

func TestFromContext(t *testing.T) {
	type testInput struct {
		ctx context.Context
	}

	type test struct {
		name   string
		input  testInput
		assert func(ctx context.Context)
	}

	tests := []test{
		{
			name: "should create a new empty context",
			input: testInput{
				ctx: context.Background(),
			},
			assert: func(ctx context.Context) {
				if _, ok := ctx.Value(contextKey).(*sync.Map); !ok {
					t.Fatal("failed to create new gocontext")
				}
			},
		},
		{
			name: "should create a new context with values",
			input: testInput{
				ctx: Add(
					FromContext(context.Background()), "test-key", "test-value",
				),
			},
			assert: func(ctx context.Context) {
				if value, ok := Get[string](ctx, "test-key"); !ok || value != "test-value" {
					t.Fatal("failed to get context key")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assert(
				FromContext(tt.input.ctx),
			)
		})
	}
}

func TestGetKeys(t *testing.T) {
	type testInput struct {
		ctx context.Context
	}

	type test struct {
		name   string
		input  testInput
		assert func(keys []string)
	}

	tests := []test{
		{
			name: "should get all context keys",
			input: testInput{
				ctx: Add(
					FromContext(context.Background()), "test-key", "test-value",
				),
			},
			assert: func(keys []string) {
				if len(keys) <= 0 {
					t.Fatal("failed to get context keys")
				}

				if !slices.Contains(keys, "test-key") {
					t.Fatal("failed to get test key from context")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assert(
				GetKeys(tt.input.ctx),
			)
		})
	}
}

func TestDelete(t *testing.T) {
	type testInput struct {
		ctx context.Context
		key string
	}

	type test struct {
		name   string
		input  testInput
		assert func(ctx context.Context, key string)
	}

	tests := []test{
		{
			name: "should delete a context value",
			input: testInput{
				ctx: Add(
					FromContext(context.Background()), "test-key", "test-value",
				),
				key: "test-key",
			},
			assert: func(ctx context.Context, key string) {
				if _, ok := Get[string](ctx, key); ok {
					t.Fatal("failed to delete context value")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Delete(tt.input.ctx, tt.input.key)
			tt.assert(tt.input.ctx, tt.input.key)
		})
	}
}
