Feature: As a passanger
    i want to be able to reserve sit  from bus of any owner 
    so that i can move one place to other.

type ticket struct {
    id
    tripId
    busId
    status
}

type Ticket interface {
    checkSeat()
    holdTicket()
}

func (t *ticket) holdTicket(tickId, tripId) {
    t.status = onhold
}

Scenario: ticket become on hold by user
    Given ticket number 12 of trip of bus 23 free
    When a user requests reservation
    Then the ticket status should be changed to "on hold",
    And the user should be redirected to checkout page  

Scenario: ticket is already reserved
    Given ticket number 12 of trip of bus 23 is reserved,     
    When a user requests reservation 
    Then the user should get response "ticket reserved"

Scenario: ticket is on hold
    Given ticket number 12 of trip of bus 23 is onhold
    When a user requests reservation
    Then the user should rseponse "ticket is onhold for 3 min, please try again in 3 minuts"

Scenario: payment success request came
    Given a user created a checkout session with paypal
    When a payment sucess callback came
    Then the ticket shoud be reserved by that user

func handleReservationCallback() {
    //const desiredShape = bodyParse(paymentMethod, body)
    
    // amount valid
    // payment gate timeout callback
    // payment gate doesn't respond in 15
    // vendor dependent implementation of the payment process
    // status update (free, reseve)
}

func timeOut()

func handleReservation(userId, ticketNO, tripId, busId) {
    // ticket struct with status free
    sit := storage.checkTicketStatus(ticketNO, tripId, busId) 
    if(sit == 'reseved) 
        ctx.json()
    if(sit == 'onhold')
        ctx.json('ticket is onhold for 3 min, please try again in 3 mintus)
    storage.holdTicket(ticket)
    return // redierct to payment
}