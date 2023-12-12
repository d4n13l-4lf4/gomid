# gomid
[![Test and coverage](https://github.com/d4n13l-4lf4/gomid/actions/workflows/test-coverage.yaml/badge.svg?branch=main)](https://github.com/d4n13l-4lf4/gomid/actions/workflows/test-coverage.yaml)
[![codecov](https://codecov.io/gh/d4n13l-4lf4/gomid/graph/badge.svg?token=JizNc5OSPg)](https://codecov.io/gh/d4n13l-4lf4/gomid)
[![pre-commit](https://img.shields.io/badge/pre--commit-enabled-brightgreen?logo=pre-commit)](https://github.com/pre-commit/pre-commit)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Reference](https://pkg.go.dev/badge/github.com/d4n13l-4lf4/gomid.svg)](https://pkg.go.dev/github.com/d4n13l-4lf4/gomid)

A middleware library written in Go

**gomid** is a middleware library which makes use of the decorator pattern to add functionality to basic functions.
Usually, it is useful to add behaviour in our application without worrying about changing our current functionality. So, this middleware library easiness this kind of changes. 

Check out our [example](https://github.com/d4n13l-4lf4/gomid-aws-example) with AWS Lambda functions to test its behaviour.

#### Usage example 
Check out this file at [main.go](https://github.com/d4n13l-4lf4/gomid/tree/main/examples/plain/main.go).
```go
// https://github.com/d4n13l-4lf4/gomid/tree/main/examples/plain/main.go

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
	fn := middleware.Wrap[func(context.Context, int) (any, error)](Adder).
		Add(Multiplier).
		Build()

	result, err := fn(context.Background(), 1)
	if err != nil {
		log.Panicln(err)
	}

	log.Printf("My result is %d", result)
}

```

#### AWS Lambda usage example

Check out a complete example at [gomid-aws-example](https://github.com/d4n13l-4lf4/gomid-aws-example).
```go
package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/d4n13l-4lf4/gomid-aws-example/auth"
	"github.com/d4n13l-4lf4/gomid-aws-example/hello"
	"github.com/d4n13l-4lf4/gomid/middleware"
)

func main() {
	greeter := hello.NewGreetingController(hello.Greet)
	chain := middleware.Wrap[func(context.Context, *events.APIGatewayProxyRequest) (any, error)](
		greeter.Greet,
	).
		Add(auth.AuthenticateUser(allowedUsers)).
		Build()

	lambda.Start(chain)
}
```
