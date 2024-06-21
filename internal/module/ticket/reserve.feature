Feature: Ticket Resevation

Scenario: user send reservation request for free ticket
    Given ticket number 12 of bus number 10 for trip of id 778 is "free"
    When user requests to reserve ticket number 12 of trip 778
    Then the ticket status should be "onhold"
    And the user should get ckeckout url 
# Scenario: user send reservation request for onhold ticket
#     Given ticket number 12 of bus number 10 is "onhold"
#     When user requests reservation 
#     Then user should get error message "ticket is onhold try later"
# Scenario: user tries to reserve already reserved ticket
#     Given ticket number 12 of bus number 10 is "reserved"
#     When user requests reservation
#     Then user should error message "ticket is already reserved"