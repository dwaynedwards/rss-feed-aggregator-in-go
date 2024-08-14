package rf

import "context"

type StateFn[T any] func(ctx context.Context, args T) (T, StateFn[T], error)

func RunStateMachine[T any](ctx context.Context, args T, start StateFn[T]) (T, error) {
	var err error
	current := start
	for {
		if ctx.Err() != nil {
			return args, ctx.Err()
		}
		args, current, err = current(ctx, args)
		if err != nil {
			return args, err
		}
		if current == nil {
			return args, nil
		}
	}
}
