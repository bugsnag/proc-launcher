Feature: POSIX signals are forwarded to the process

    Background:
        Given I build the executable
        When I run the executable with arguments:
            | sleep | 30 |

    Scenario Outline: Termination with signals
        Given I am using a POSIX-compliant system
        When I send <sig> to the executable
        Then the process exited with signal <sig>

        Examples:
            | sig     |
            | SIGABRT |
            | SIGALRM |
            | SIGTERM |
            | SIGQUIT |
            | SIGINT  |
            | SIGHUP  |
