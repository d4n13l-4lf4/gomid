package middleware_test

import (
	"context"
	"errors"
	"testing"

	"github.com/cucumber/godog"
	"github.com/d4n13l-4lf4/gomid/internal"
	"github.com/d4n13l-4lf4/gomid/middleware"
	"github.com/stretchr/testify/assert"
)

const (
	FunctionTypeError = "ERR_FUNCTION_TYPE"
	AllowedArgsError  = "ERR_ALLOWED_ARGS"
	FirstInArgError   = "ERR_FIRST_IN_ARG"
	LastOutArgError   = "ERR_LAST_OUT_ARG"
	Zero              = 0
)

type (
	testChainKey   struct{}
	testChainError struct {
		chain  *middleware.Wrapper
		baseFn any
	}
	testChainSuccess struct {
		chain *middleware.Wrapper
		out   any
		err   error
	}
)

func useWrongBaseFunction(ctx context.Context, wrappedType string) (context.Context, error) {
	wrappedFn := getWrappedFunction(wrappedType)
	chain := middleware.Wrap(wrappedFn)
	test := &testChainError{
		chain:  chain,
		baseFn: wrappedFn,
	}

	return context.WithValue(ctx, &testChainKey{}, test), nil
}

func shouldPanicAndShowErrorMessage(ctx context.Context, msg string) error {
	var asserter internal.Asserter
	test := ctx.Value(&testChainKey{}).(*testChainError)

	assert.PanicsWithValue(&asserter, msg, func() {
		test.chain.Build()
	})

	return asserter.Error()
}

func wantToAddMiddlewares(ctx context.Context) (context.Context, error) {
	wrappedFn := func(_ context.Context, a int) (int, error) {
		return a + 1, nil
	}
	chain := middleware.Wrap(wrappedFn)
	test := &testChainSuccess{
		chain: chain,
	}

	return context.WithValue(ctx, &testChainKey{}, test), nil
}

func haveBasefunctionWithAnError(ctx context.Context, errMsg string) (context.Context, error) {
	wrappedFn := func(_ context.Context, a int) (int, error) {
		return 0, errors.New(errMsg)
	}
	chain := middleware.Wrap(wrappedFn)
	test := &testChainSuccess{
		chain: chain,
	}

	return context.WithValue(ctx, &testChainKey{}, test), nil
}

func addMiddlewares(ctx context.Context) (context.Context, error) {
	test := ctx.Value(&testChainKey{}).(*testChainSuccess)
	adder := func(nxt middleware.Next) middleware.Next {
		return func(_ context.Context, a any) (any, error) {
			numA := a.(int)

			return nxt(ctx, numA+1)
		}
	}
	multiplier := func(nxt middleware.Next) middleware.Next {
		return func(_ context.Context, a any) (any, error) {
			numA := a.(int)

			return nxt(ctx, numA*2)
		}
	}

	chain := test.chain.
		Add(adder).
		Add(multiplier)

	test.chain = chain

	return context.WithValue(ctx, &testChainKey{}, test), nil
}

func addMiddlareWithError(ctx context.Context, errMsg string) (context.Context, error) {
	test := ctx.Value(&testChainKey{}).(*testChainSuccess)
	triggerError := func(nxt middleware.Next) middleware.Next {
		return func(_ context.Context, a any) (any, error) {
			return nil, errors.New(errMsg)
		}
	}

	chain := test.chain.
		Add(triggerError)

	test.chain = chain

	return context.WithValue(ctx, &testChainKey{}, test), nil
}

func executeMiddleware(ctx context.Context) (context.Context, error) {
	test := ctx.Value(&testChainKey{}).(*testChainSuccess)
	chain := test.chain
	chainedFn := chain.Build().(func(context.Context, int) (any, error))

	//nolint:contextcheck
	test.out, test.err = chainedFn(context.TODO(), 1)

	return context.WithValue(ctx, &testChainKey{}, test), nil
}

func shouldGetMyNumber(ctx context.Context, num int) error {
	var asserter internal.Asserter

	assertions := assert.New(&asserter)
	test := ctx.Value(&testChainKey{}).(*testChainSuccess)

	assertions.Equal(num, test.out)
	assertions.NoError(test.err)

	return asserter.Error()
}

func shouldGetError(ctx context.Context, errMsg string) error {
	var asserter internal.Asserter

	assertions := assert.New(&asserter)
	test := ctx.Value(&testChainKey{}).(*testChainSuccess)

	assertions.Zero(test.out)
	assertions.EqualError(test.err, errMsg)

	return asserter.Error()
}

func getWrappedFunction(functionType string) any {
	switch functionType {
	case FunctionTypeError:
		return 1
	case AllowedArgsError:
		return func() {}
	case FirstInArgError:
		return func(_ int, _ int) (int, error) {
			return Zero, nil
		}
	case LastOutArgError:
		return func(_ context.Context, a int) (int, int) { return Zero, Zero }
	default:
		return func() {}
	}
}

func InitializeChainScenario(sc *godog.ScenarioContext) {
	sc.Given("^I want to build a middleware chain$", func() {})
	sc.When("^I use a ([A-Z_]+) function$", useWrongBaseFunction)
	sc.Then("^It should panic and show an error ([A-Za-z\\s\\.\\*]+)$", shouldPanicAndShowErrorMessage)

	sc.Given("^I want to add middlewares to my base function$", wantToAddMiddlewares)
	sc.When("^I add middlewares to my base function$", addMiddlewares)
	sc.Step("^I add a middleware with an error ([a-z\\s]+)$", addMiddlareWithError)
	sc.Step("^I execute my functionality$", executeMiddleware)

	sc.Then("^I should get a number (\\d+)$", shouldGetMyNumber)
	sc.Then("^I should get an error ([a-z\\s]+)$", shouldGetError)

	sc.Given("^I have a base function which returns an error ([a-z\\s]+)$", haveBasefunctionWithAnError)
}

func TestChainFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeChainScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
