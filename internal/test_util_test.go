package internal_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/d4n13l-4lf4/gomid/internal"
	"github.com/stretchr/testify/assert"
)

func TestAsserter(t *testing.T) {
	tests := []struct {
		Name string
	}{
		{
			"it should return a formatted error when called",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			assertions := assert.New(t)
			asserter := internal.NewAsserter()
			errTest := errors.New("test error")

			expectedErr := fmt.Sprintf("error %s", errTest.Error())
			asserter.Errorf("error %s", errTest.Error())

			assertions.EqualError(asserter.Error(), expectedErr)
		})
	}
}
