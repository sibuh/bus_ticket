Feature: Ticket Resevation

  Scenario: user send reservation request for free ticket
    Given a "Free" ticket
    When user requests to reserve ticket
    Then the ticket status should be "Onhold"
    And checkout session request should be sent

  Scenario: checkout session create request responds success
   Given a "Free" ticket
   When create checkout session succeeds for reserving ticket request
   Then checkout session should be stored
   And onhold-timeout process should be scheduled
   And the user should get checkout url

  Scenario: checkout session create request fails
    Given a "Free" ticket
    When checkout session creation fails during reserve ticket request
    Then user should get error message "failed to create checkout session"
    And the ticket status should be "Free"

  Scenario: user tries to reserve already held ticket
    Given a "Onhold" ticket
    When user requests to reserve ticket
    Then user should get error message "ticket is onhold please try later"

  Scenario: user tries to reserve already reserved ticket
    Given a "Reserved" ticket
    When user requests to reserve ticket
    Then user should get error message "ticket is already reserved please try to reserve free ticket"

