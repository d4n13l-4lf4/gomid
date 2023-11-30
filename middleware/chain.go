package middleware

import (
	"context"
	"log"
	"reflect"
)

const allowedArgs = 2

type (
	Wrapper struct {
		wrapped     any
		middlewares []Middleware
	}

	Middleware func(nxt Next) Next

	Next func(context.Context, any) (any, error)
)

func Wrap(fn any) *Wrapper {
	return &Wrapper{
		wrapped: fn,
	}
}

func (w *Wrapper) Add(midd Middleware) *Wrapper {
	w.middlewares = append(w.middlewares, midd)

	return w
}

func (w *Wrapper) verifyWrapped(val reflect.Value) {
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

func (w *Wrapper) Build() any {
	wrappedVal := reflect.ValueOf(w.wrapped)

	w.verifyWrapped(wrappedVal)

	head := func(wrapped reflect.Value) Next {
		return func(ctx context.Context, data any) (any, error) {
			in := []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(data)}
			out := wrapped.Call(in)
			if err, ok := out[1].Interface().(error); ok {
				return out[0].Interface(), err
			}

			return out[0].Interface(), nil
		}
	}(wrappedVal)

	for i := len(w.middlewares) - 1; i >= 0; i-- {
		next := w.middlewares[i](head)
		head = next
	}

	wrappedType := wrappedVal.Type()
	decoratedFn := reflect.ValueOf(head)

	inner := func(decorated reflect.Value) func(args []reflect.Value) (results []reflect.Value) {
		return func(args []reflect.Value) (results []reflect.Value) {
			return decorated.Call(args)
		}
	}(decoratedFn)

	interfaceType := reflect.TypeOf((*any)(nil)).Elem()

	in := []reflect.Type{wrappedType.In(0), wrappedType.In(1)}
	out := []reflect.Type{interfaceType, wrappedType.Out(1)}
	genericFn := reflect.FuncOf(in, out, false)

	outFn := reflect.MakeFunc(genericFn, inner)

	return outFn.Interface()
}
