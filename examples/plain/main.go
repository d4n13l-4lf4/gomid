package main

import (
	"context"
	"log"

	"github.com/d4n13l-4lf4/gomid/middleware"
)

func Adder(ctx context.Context, a int) (int, error) {
	log.Println(ctx)
	return a + 1, nil
}

func Multiplier(next middleware.Next) middleware.Next {
	return func(ctx context.Context, a any) (any, error) {
		newA := a.(int)

		return next(ctx, newA*2)
	}
}

func main() {
	fn := middleware.Wrap(Adder).
		Add(Multiplier).
		Build().(func(context.Context, int) (any, error))

	result, err := fn(context.Background(), 1)
	if err != nil {
		log.Panicln(err)
	}

	log.Printf("My result is %d", result)
}
