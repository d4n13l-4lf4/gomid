package middleware

import (
	"context"
	"log"
	"reflect"
)

const allowedArgs = 2

type (
	// Wrapper wraps a base function to be decorated with middlewares.
	Wrapper[R any] struct {
		wrapped     any
		middlewares []Middleware
	}

	// Middleware is a function which receives the next function in the chain.
	Middleware func(nxt Next) Next

	// Next is a function which performs actual processing of input and output.
	Next func(context.Context, any) (any, error)
)

// Wrap creates a ready to use wrapper for adding middlewares.
func Wrap[R any](fn any) *Wrapper[R] {
	return &Wrapper[R]{
		wrapped: fn,
	}
}

// Add adds a new middleware to be used in the chain.
func (w *Wrapper[R]) Add(midd Middleware) *Wrapper[R] {
	w.middlewares = append(w.middlewares, midd)

	return w
}

func (w *Wrapper[R]) verifyWrapped(val reflect.Value) {
	wrappedType := val.Type()

	if wrappedType.Kind() != reflect.Func {
		log.Panicf("Type %s is not function.\n", wrappedType.String())
	}

	inArgs := wrappedType.NumIn()
	outArgs := wrappedType.NumOut()

	if inArgs != outArgs || inArgs != allowedArgs {
		log.Panicf("Function does not conform with signature %s.\n", reflect.TypeOf((*Next)(nil)))
	}

	ctxType := reflect.TypeOf((*context.Context)(nil)).Elem()
	errType := reflect.TypeOf((*error)(nil)).Elem()
	firstInArg := wrappedType.In(0)
	lastOutArg := wrappedType.Out(1)

	if !firstInArg.Implements(ctxType) {
		log.Panicf("First argument %s does not implements context.Context interface.\n", firstInArg.String())
	}

	if !lastOutArg.Implements(errType) {
		log.Panicf("Last output argument %s does not implements error interface.\n", lastOutArg.String())
	}
}

// Build builds a middleware chain with defined middlewares.
func (w *Wrapper[R]) Build() R {
	wrappedVal := reflect.ValueOf(w.wrapped)

	w.verifyWrapped(wrappedVal)

	head := func(wrapped reflect.Value) Next {
		return func(ctx context.Context, data any) (any, error) {
			in := []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(data)}
			out := wrapped.Call(in)
			firstOut := out[0].Interface()
			if err, ok := out[1].Interface().(error); ok {
				return firstOut, err
			}

			return firstOut, nil
		}
	}(wrappedVal)

	for i := len(w.middlewares) - 1; i >= 0; i-- {
		next := w.middlewares[i](head)
		head = next
	}

	wrappedType := wrappedVal.Type()
	chainedFn := reflect.ValueOf(head)

	inner := func(chained reflect.Value) func(args []reflect.Value) (results []reflect.Value) {
		return func(args []reflect.Value) (results []reflect.Value) {
			return chained.Call(args)
		}
	}(chainedFn)

	interfaceType := reflect.TypeOf((*any)(nil)).Elem()

	in := []reflect.Type{wrappedType.In(0), wrappedType.In(1)}
	out := []reflect.Type{interfaceType, wrappedType.Out(1)}
	genericFn := reflect.FuncOf(in, out, false)

	outFn := reflect.MakeFunc(genericFn, inner)

	return outFn.Interface().(R)
}
