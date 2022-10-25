package values

import (
	"context"
	"reflect"

	"github.com/pkg/errors"
)

func Zero[T any]() T {
	var zero T
	return zero
}

func IsZero[T any](value T) bool {
	return reflect.DeepEqual(value, Zero[T]())
}

func ContextAssertion[T any](ctx context.Context, key string) (T, error) {
	v := ctx.Value(key)
	x, ok := v.(T)
	if !ok {
		return Zero[T](), errors.Errorf("assertion error, values mismatch: %T", v)
	}
	return x, nil
}
