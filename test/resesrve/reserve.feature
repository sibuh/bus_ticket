Feature: As a passanger
    i want to able to reserve sit  from bus of any owner 
    so that i can move one place to other.


Background: 
    Given there is unreserved ticket
    And there is ticket which is onhold

Scenario: user reserves ticket successfully
    When user trys to reserve "free" ticket
    Then user should get ticket reserved. 
    And user should get reservation receipt

Scenario: user fails to reserve ticket
    When user trys to reserve "onhold"ticket
    Then user should get "<error>" response

Examples:
|error            |
|ticket is on hold|

Scenario: two users try to reserve single sit concurrently
    When two users try to reserve one sit at a time 
    Then one user should get ticket reserved 
    And the other should get "<error>"
Examples:
|error            |
|ticket is on hold|









    