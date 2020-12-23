Feature: The contents of standard output streams are relayed to/from the process

    Background:
        Given I build the executable

    Scenario: Forwarding process output to stdout
        When I run the executable with arguments:
            | bash | -c | echo Hello |
        Then "Hello" is present in the standard output stream

    Scenario: Forward process error output to stderr
        When I run the executable with arguments:
            | bash | -c | echo Goodbye >&2 |
        Then "Goodbye" is present in the standard error stream
