Feature: Start server
    Scenario: I connect to server control
        Given ready minecraft server

    Scenario: I can set a config
        When I set the config to
        """
        {
          "hardcore": true,
          "max-players": 55,
          "motd": "A server"
        }
        """
        Then I have no errors
        When the config equal to
        """
        {
          "hardcore": true,
          "max-players": 55,
          "motd": "A server"
        }
        """
        Then I have no errors
        