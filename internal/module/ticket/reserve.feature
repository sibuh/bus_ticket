Feature: Ticket Resevation

  Scenario: user send reservation request for free ticket
    Given ticket number 12 of bus number 10 for trip of id 778 is "Free"
    When user requests to reserve ticket number 12 of trip 778
    Then the ticket status should be "Onhold"
    And the user should get checkout url

 Scenario: user send reservation request for free ticket but create checkout session fails
    Given ticket number 12 of bus number 10 for trip of id 779 is "Free"
    When user requests to reserve ticket number 12 of trip 778
    Then the ticket status should be "Free"
    And user should get error message "failed to create checkout session"

  Scenario: user tries to reserve already held ticket
    Given ticket number 12 of bus number 10 for trip of id 778 is "Onhold"
    When user requests to reserve ticket number 12 of trip 778
    Then user should get error message "ticket is onhold please try later"

  Scenario: user tries to reserve already reserved ticket
    Given ticket number 12 of bus number 10 for trip of id 778 is "Reserved"
    When user requests to reserve ticket number 12 of trip 778
    Then user should get error message "ticket is already reserved please try to reserve free ticket"

  Scenario: checkout session timeout
    Given checkout session is created
      | id | ticket number | bus number | time     |
      |  1 |            12 |   10       | 12-12-20 |
    When ticket reservstion do not succeed with in 10s duration
    Then check payment status on payment gateway

  Scenario: payment gateway returned successful purchase status
    Given checkout session is created
      | id | ticket number | bus number | time     |
      |  1 |            12 |   10       | 12-12-20 |
    And payment status is requested for checkout session
    When payment status checkout session returns "successful purchase status" for checkout session
    Then ticket must be set to "Reserved" status

  Scenario: payment gateway returned pending purchase status
    Given checkout session is created
      | id | ticket number | bus number | time     |
      |  1 |            12 |   10       | 12-12-20 |
    And payment status is requested for checkout session
    When payment status for checkout session returns "pending"
    Then cancel checkout session is sent to payment gateway

  Scenario: payment cancelation successful
    When payment cancelation response is successful
    Then ticket must be set "Free" for sale

  Scenario: payment cancelation not successful
  # retry
