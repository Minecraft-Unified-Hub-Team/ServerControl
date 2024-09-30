Feature: Stop server
    Scenario: Server has "alive" state
        When I get server state "alive"
        Then I have no errors

    Scenario: I can stop server
        When I stop server
        Then I have no errors

    Scenario: Server has "stopped" state
        When I get server state "stopped"
        Then I have no errors
