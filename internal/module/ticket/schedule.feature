Feature: schedule ontimeout process

Scenario: send check payment status request to gateway
Given checkout session is successfully created 
When payment do not complete within 10 seconds
Then payment status check request should be sent to gateway


