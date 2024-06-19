Feature: Ticket Resevation

Scenario: Hold ticket
    Given ticket number 12 of bus number 10 is "free"
    When user requests reservation 
    Then the ticket status should be "onhold"
    And the user should get ckeckout url 
# Scenario: ticket onhold
#     Given ticket number 12 of bus number 10 is "onhold"
#     When user requests reservation 
#     Then user should get error message "ticket is onhold try later"
# Scenario: already reserved ticket
#     Given ticket number 12 of bus number 10 is "reserved"
#     When user requests reservation
#     Then user should error message "ticket is already reserved"