package internal

import (
	"fmt"
)

type (
	// Asserter is an assertion implementation for tests.
	Asserter struct {
		err error
	}
)

// NewAsserter creates a new *Asserter.
func NewAsserter() *Asserter {
	return &Asserter{}
}

// Errorf formats an error according to fmt.Errorf rules.
func (a *Asserter) Errorf(format string, args ...any) {
	a.err = fmt.Errorf(format, args...)
}

// Error returns the inner error.
func (a *Asserter) Error() error {
	return a.err
}
