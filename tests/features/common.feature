Feature: Start server
    Scenario: I connect to server control
        Given ServerControl is up
        And I have no errors
        When I connect to service control
        Then I have no errors

    Scenario: Ping works
        When I ping to the server
        Then I have no errors

    Scenario: I can install server
        When I install "1.20.6-50.1.3" server version
        Then I have no errors

    Scenario: Server has "Stopped" state
        When I get server state 
        """
        {
          "State": "Stopped"
        }
        """
        Then I have no errors

    Scenario: I can start server
        When I start server
        Then I have no errors

    Scenario: Server has "Alive" state
        When I get server state 
        """
        {
          "State": "Alive"
        }
        """
        Then I have no errors

    Scenario: I can stop server
        When I stop server
        Then I have no errors

    Scenario: Server has "Stopped" state
        When I get server state 
        """
        {
          "State": "Stopped"
        }
        """
        Then I have no errors

    Scenario: I can start server
        When I start server
        Then I have no errors

    Scenario: Server has "Alive" state
        When I get server state 
        """
        {
          "State": "Alive"
        }
        """
        Then I have no errors

    Scenario: I can uninstall server
        When I uninstall server
        Then I have no errors
