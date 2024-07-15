Feature: Ticket reservation

 Scenario: checkout session create request responds success
   Given a free ticket
    When create checkout session succeeds for reserving ticket request
    Then checkout session should be stored
    And the user should get checkout url