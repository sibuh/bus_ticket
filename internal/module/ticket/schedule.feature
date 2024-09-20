Feature: schedule ontimeout process

  Scenario:
    When payment status check request is scheduled for checkout session
    Then no payment status check request should be sent within 2s

  Scenario:
    When payment status check request is scheduled for checkout session
    Then payment status check request should be sent to payment gateway after 2s

  Scenario:
    Given payment status check request is scheduled for checkout session
    When success or failure callback arrives for checkout session
    Then scheduled process should be terminated
