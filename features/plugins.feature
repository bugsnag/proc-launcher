Feature: The launcher can have plugins respond to events

    Background:
        Given I build the extension executable

    Scenario: Custom handling for stdout
        When I run the extension executable with arguments:
            | bash | -c | echo Hello |
        Then "you said: Hello" is present in the standard output stream

    Scenario: Custom handling for stderr
        When I run the extension executable with arguments:
            | bash | -c | echo Hello >&2 |
        Then "something bad happened: Hello" is present in the standard output stream

    Scenario: Custom handling at exit
        When I run the extension executable with arguments:
            | bash | -c | exit 47 |
        Then "process terminated, code 47" is present in the standard output stream
