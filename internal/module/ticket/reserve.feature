Feature: Ticket Resevation

Scenario: user send reservation request for free ticket
    Given ticket number 12 of bus number 10 for trip of id 778 is "Free"
    When user requests to reserve ticket number 12 of trip 778
    Then the ticket status should be "Onhold"
    And the user should get checkout url 
Scenario: user tries to reserve already held ticket
    Given ticket number 12 of bus number 10 for trip of id 778 is "Onhold"
    When user requests to reserve ticket number 12 of trip 778
    Then user should get error message "ticket is onhold please try later"
Scenario: user tries to reserve already reserved ticket
    Given ticket number 12 of bus number 10 for trip of id 778 is "Reserved"
    When user requests to reserve ticket number 12 of trip 778
    Then user should get error message "ticket is already reserved please try to reserve free ticket"
Scenario: held ticket released to free incase payment delays
    Given ticket number 12 of bus number 10 for trip of id 778 is "Onhold"
    When 
