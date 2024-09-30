Feature: Start server
    Scenario: I connect to server control
        When I connect to service control
        Then I have no errors

    Scenario: Ping works
        And I ping to the server
        Then I have no errors

    Scenario: I can install server
        When I install "1.20.6-50.1.3" server version
        Then I have no errors

    Scenario: I can start server
        When I start server
        Then I have no errors

    Scenario: Server has "alive" state
        When I get server state "alive"
        Then I have no errors
