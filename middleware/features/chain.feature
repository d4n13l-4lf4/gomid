Feature: Building a middleware chain
    In order to add functionality to a base function
    As a Gopher
    I need to be able to use the decorator pattern.

    Scenario: A middleware chain returns my number successfully
        Given I want to add middlewares to my base function
        When I add middlewares to my base function
        And I execute my functionality
        Then I should get a number <number>

        Examples:
        | number |
        | 5 |
    
    Scenario: My base function returns an error successfully
        Given I have a base function which returns an error <error>
        When I add middlewares to my base function
        And I execute my functionality
        Then I should get an error <error>

        Examples:
        | error |
        | base function error |
    
    Scenario: A middleware chain returns an error successfully
        Given I want to add middlewares to my base function
        When I add middlewares to my base function
        And I add a middleware with an error <error>
        And I execute my functionality
        Then I should get an error <error>

        Examples:
        | error |
        | could not achieve what you wanted |
    
    
    Scenario: A middleware chain cannot be built successfully
        Given I want to build a middleware chain
        When I use a <base> function
        Then It should panic and show an error <message>

        Examples:
        | base | message |
        | ERR_FUNCTION_TYPE | Type int is not function.\n |
        | ERR_ALLOWED_ARGS  | Function does not conform with signature *middleware.Next.\n |
        | ERR_FIRST_IN_ARG | First argument int does not implements context.Context interface.\n |
        | ERR_LAST_OUT_ARG | Last output argument int does not implements error interface.\n |