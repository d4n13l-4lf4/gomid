package middleware_test

import (
	"context"
	"testing"

	"github.com/cucumber/godog"
	"github.com/d4n13l-4lf4/gomid/middleware"
)

func wrapAGenericFunction(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func runWrapped() error {
	middleware.Middleware()

	return nil
}

func shouldBeCalledInTheMidst() error {
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Given(`^I wrap a function with a middleware$`, wrapAGenericFunction)
	ctx.When(`^I run the wrapped function$`, runWrapped)
	ctx.Then(`^it should be called in the midst of middleware execution$`, shouldBeCalledInTheMidst)
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
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
