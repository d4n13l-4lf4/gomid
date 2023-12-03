# gomid
[![codecov](https://codecov.io/gh/d4n13l-4lf4/gomid/graph/badge.svg?token=JizNc5OSPg)](https://codecov.io/gh/d4n13l-4lf4/gomid)
[![pre-commit](https://img.shields.io/badge/pre--commit-enabled-brightgreen?logo=pre-commit)](https://github.com/pre-commit/pre-commit)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A middleware library written in Go

#### Usage example 
```go
// https://github.com/d4n13l-4lf4/gomid/examples/plain/main.go

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

```