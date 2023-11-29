Feature: middleware functions
In order to add functionality independently
As a Gopher
I need to be able to decorate functions with middlewares

    Scenario: Add a logging middleware to a base function
        Given I wrap a function with a middleware
        When I run the wrapped function
        Then it should be called in the midst of middleware execution
    