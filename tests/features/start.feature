Feature: Start server
    Scenario: I connect to server control
        When I connect to service control
        Then I have no errors

    Scenario: Ping works
        And I ping to the server
        Then I have no errors

    Scenario: I can start server
        When I start server
        Then I have no errors