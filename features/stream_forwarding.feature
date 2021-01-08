Feature: The contents of standard output streams are relayed to/from the process

    Background:
        Given I build the executable

    Scenario: Forwarding process output to stdout
        When I run the executable with arguments:
            | bash | -c | echo Hello |
        Then "Hello" is present in the standard output stream

    Scenario: Forwarding intermittent process output to stdout
        When I run the executable with arguments:
            | bash | -c | echo "Hello There\nfriend" |
        Then "Hello There\nfriend" is present in the standard output stream

    Scenario: Forward process error output to stderr
        When I run the executable with arguments:
            | bash | -c | echo Goodbye >&2 |
        Then "Goodbye" is present in the standard error stream

    Scenario: Forward intermittent process error output to stderr
        When I run the executable with arguments:
            | bash | -c | echo "Goodbye,\nfriend" >&2 |
        Then "Goodbye,\nfriend" is present in the standard error stream

    Scenario: Forward delayed process output
        When I run the executable with arguments:
            | bash | -c | echo "Hello" && sleep 1 && echo "Goodbye,\nfriend" |
        Then "Hello\nGoodbye,\nfriend" is present in the standard output stream
