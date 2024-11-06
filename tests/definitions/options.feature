Feature: Start server
    Scenario: I can change an option
        Given option "default_timeout" equal to "10"
        When I set the option "default_timeout" with value "15"
        Then I have no errors
        And option "default_timeout" equal to "15"
    
    Scenario: I can't change an unknown option
        When I set the option "my_variable_32" with value "15"
        Then I have the error "there isn't an option with the name my_variable_32"